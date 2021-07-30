//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/07/2021    CLH             Initial version
// 07/29/2021    CLH             Add SSUrl to CommonInputs

package tkesdk

import (
	"encoding/hex"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
	"github.com/IBM/ibm-hpcs-tke-sdk/ep11cmds"
)

/*----------------------------------------------------------------------------*/
/* Structures for holding input and output values for TKE SDK functions       */
/*----------------------------------------------------------------------------*/

// Structure containing information on an installed administrator
type ReturnedAdminInfo struct {
	AdminName string
	AdminSKI  string
}

// Structure containing information describing a crypto unit assigned to
// the service instance
type HsmInfo struct {
	HsmId               string
	HsmLocation         string
	HsmType             string
	SignatureThreshold  int
	RevocationThreshold int
	Admins              []ReturnedAdminInfo
	NewMKStatus         string
	NewMKVP             string
	CurrentMKStatus     string
	CurrentMKVP         string
}

// Structure describing administrators to be created or used
type AdminInfo struct {
	Name string
	Key  string
		// This identifies the administrator signature key to be used.
		// For initial development, this will be the fully qualified path
		// and file name of a signature key file.
		// When user-defined signing services are supported, the signing
		// service will define how this field is set.
	Token string
		// Credential giving access to the administrator signature key.
		// For initial development, this will be the file password.
		// When user-defined signing services are supported, the signing
		// service will define how this field is set.
}

// Structure representing the hsm_config section of a resource block
type HsmConfig struct {
	SignatureThreshold  int
	RevocationThreshold int
	Admins              []AdminInfo
}

/*----------------------------------------------------------------------------*/
/* Collects and returns information on how the crypto units assigned to a     */
/* service instance are configured.                                           */
/*----------------------------------------------------------------------------*/
func Query(ci common.CommonInputs) ([]HsmInfo, error) {
	hsmInfo, _, err := internalQuery(ci)
	return hsmInfo, err
}

/*----------------------------------------------------------------------------*/
/* Function used internally to query the crypto unit configuration.           */
/*                                                                            */
/* Returns additional information used by other TKE SDK functions.            */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/*                                                                            */
/* Outputs:                                                                   */
/* []HsmInfo -- an array of structures with the current configuration         */
/*      settings for each crypto unit in the service instance                 */
/* []common.DomainEntry -- identifies the set of crypto units assigned to     */
/*      the service instance                                                  */
/* error -- reports an error encountered when running the function, nil if    */
/*      no error found                                                        */
/*----------------------------------------------------------------------------*/
func internalQuery(ci common.CommonInputs) ([]HsmInfo, []common.DomainEntry, error) {

	// Create an empty output array
	hsmInfo := make([]HsmInfo, 0)

	// Create an empty domains array
	domains := make([]common.DomainEntry, 0)

	// Query to see what crypto units are assigned to the service instance
	domains, err := getDomains(ci)
	if err != nil {
		return hsmInfo, domains, err
	}

	for _, domain := range domains {

		// Create an empty structure for this domain
		nextHsm := HsmInfo{}

		nextHsm.HsmId = domain.Hsm_id
		nextHsm.HsmLocation = domain.Location
		nextHsm.HsmType = domain.Type

		// Query the signature thresholds
		domAttr, _, err := ep11cmds.QueryDomainAttributes(ci, domain)
		if err != nil {
			return hsmInfo, domains, err
		}
		nextHsm.SignatureThreshold = int(domAttr.SignatureThreshold)
		nextHsm.RevocationThreshold = int(domAttr.RevocationSignatureThreshold)

		// Query domain administrators
		domAdminSKIs, err := ep11cmds.QueryDomainAdmins(ci, domain)
		if err != nil {
			return hsmInfo, domains, err
		}
		nextHsm.Admins = make([]ReturnedAdminInfo, len(domAdminSKIs))
		for j := 0; j < len(domAdminSKIs); j++ {
			nextHsm.Admins[j].AdminSKI = hex.EncodeToString(domAdminSKIs[j])

			name, err := ep11cmds.QueryDomainAdminName(ci, domain, domAdminSKIs[j])
			if err != nil {
				return hsmInfo, domains, err
			}
			nextHsm.Admins[j].AdminName = name
		}

		// Query master key register state and verification pattern
		domainInfo, err := ep11cmds.QueryDomainInfo(ci, domain)
		if err != nil {
			return hsmInfo, domains, err
		}

		nextHsm.NewMKStatus = convertMKStatusToString(domainInfo.NewMKStatus)
		nextHsm.NewMKVP = hex.EncodeToString(domainInfo.NewMKVP)

		nextHsm.CurrentMKStatus = convertMKStatusToString(domainInfo.CurrentMKStatus)
		nextHsm.CurrentMKVP = hex.EncodeToString(domainInfo.CurrentMKVP)

		hsmInfo = append(hsmInfo, nextHsm)
	}

	return hsmInfo, domains, nil
}

/*----------------------------------------------------------------------------*/
/* Returns an appropriate string for the master key status value              */
/*----------------------------------------------------------------------------*/
func convertMKStatusToString(status int) string {
	if status == ep11cmds.MK_STATUS_EMPTY {
		return "Empty"
	} else if status == ep11cmds.CMK_STATUS_VALID {
		return "Valid"
	} else if status == ep11cmds.NMK_STATUS_FULL_UNCOMMITTED {
		return "Full Uncommitted"
	} else if status == ep11cmds.NMK_STATUS_FULL_COMMITTED {
		return "Full Committed"
	} else {
		return ""
	}
}
