//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 05/04/2021    CLH             Adapt for TKE SDK
// 07/30/2021    CLH             Add SSUrl to CommonInputs

package ep11cmds

import (
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Finalizes the pending wrapping key register                                */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain whose pending wrapping key register   */
/*    is to be finalized                                                      */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func FinalizeWK(ci common.CommonInputs, de common.DomainEntry,
	sigkeys []string, sigkeySkis []string, sigkeyTokens []string) error {

	// Get the verification pattern for the pending wrapping key register
	domainInfo, err := QueryDomainInfo(ci, de)
	if err != nil {
		return err
	}

	htpRequestString, err := FinalizeWKReq(ci, de, domainInfo.NewMKVP,
		sigkeys, sigkeySkis, sigkeyTokens)
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
/* Creates the HTPRequest for finalizing the pending wrapping key register    */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain whose pending wrapping key register   */
/*    is to be finalized                                                      */
/* []byte -- the verification pattern of the pending wrapping key register    */
/*    to be finalized                                                         */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the HTPRequest string with the signed CPRB for the command       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func FinalizeWKReq(ci common.CommonInputs, de common.DomainEntry,
	vp []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_FINALIZE_WK
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	adminBlk.CmdInput = vp
	return CreateSignedHTPRequest(ci, de, adminBlk, sigkeys, sigkeySkis, sigkeyTokens)
}
