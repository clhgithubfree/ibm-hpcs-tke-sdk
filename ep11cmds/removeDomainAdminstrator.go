//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/29/2021    CLH             Adapt for TKE SDK
// 07/30/2021    CLH             Add SSUrl to CommonInputs

package ep11cmds

import (
	"encoding/hex"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Removes an administrator                                                   */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain with the administrator to be removed  */
/* string -- the Subject Key Identifier of the administator to be removed     */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func RemoveDomainAdministrator(ci common.CommonInputs, de common.DomainEntry,
	ski string, sigkeys []string, sigkeySkis []string, sigkeyTokens []string) error {

	// Convert from hexadecimal string to []byte
	skibytes, err := hex.DecodeString(ski)
	if err != nil {
		return err
	}

	htpRequestString, err := RemoveDomainAdminReq(ci, de, skibytes, sigkeys,
		sigkeySkis, sigkeyTokens)
	if err != nil {
		return err
	}

	req, err := common.CreatePostHsmsRequest(ci, de.Hsm_id, htpRequestString)
	if err != nil {
		return err
	}

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return err
	}

	_, err = buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return err
	}

	return nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for removing a domain administrator                 */
/*----------------------------------------------------------------------------*/
func RemoveDomainAdminReq(ci common.CommonInputs, de common.DomainEntry,
	ski []byte, sigkeys []string, sigkeySkis []string, sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_DOM_ADMIN_LOGOUT
	// DomainID, ModuleID, and TransactionCounter get filled in later when sending the request
	adminBlk.CmdInput = ski
	return CreateSignedHTPRequest(ci, de, adminBlk, sigkeys, sigkeySkis, sigkeyTokens)
}
