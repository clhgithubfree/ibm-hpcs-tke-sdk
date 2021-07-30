//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/07/2021    CLH             Adapt for TKE SDK
// 07/30/2021    CLH             Add SSUrl to CommonInputs

package ep11cmds

import (
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Queries the domain control points                                          */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain to be queried                         */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the domain control points (16 bytes long)                        */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDomainControlPoints(ci common.CommonInputs, de common.DomainEntry) ([]byte, error) {

	htpRequestString := QueryDomainControlPointsReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req, err := common.CreatePostHsmsRequest(ci, de.Hsm_id, htpRequestString)
	if err != nil {
		return nil, err
	}

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return nil, err
	}

	adminRspBlk, err := buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return nil, err
	}

	return adminRspBlk.CmdOutput, nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for querying domain control points                  */
/*----------------------------------------------------------------------------*/
func QueryDomainControlPointsReq(cryptoModuleIndex int, domainIndex int) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DOM_CTRLPOINTS
	adminBlk.DomainID = BuildAdminDomainIndex(domainIndex)
	// module ID not used for queries
	// transaction counter not used for queries
	// no input parameters
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}
