//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/09/2021    CLH             Adapt for TKE SDK

package ep11cmds

import (
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Loads a random value in one of the wrapping key registers.                 */
/*                                                                            */
/* If the current wrapping key register is empty, it is loaded with a random  */
/* value.                                                                     */
/*                                                                            */
/* If the current wrapping key register is not empty but the pending wrapping */
/* key register is empty, the pending wrapping key register is loaded with a  */
/* random value.                                                              */
/*                                                                            */
/* If both the current wrapping key register and pending wrapping key         */
/* register are not empty, an error is returned.                              */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain where a random value is to be loaded  */
/*    in one of the wrapping key registers                                    */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/* []byte -- the verification pattern of the generated master key value       */
/*----------------------------------------------------------------------------*/
func CreateRandomWK(ci common.CommonInputs, de common.DomainEntry,
	sigkeys []string, sigkeySkis []string, sigkeyTokens []string) (error, []byte) {

	htpRequestString, err := CreateRandomWKReq(ci, de, sigkeys, sigkeySkis, sigkeyTokens)
	if err != nil {
		return err, nil
	}

	req, err := common.CreatePostHsmsRequest(ci, de.Hsm_id, htpRequestString)
	if err != nil {
		return err, nil
	}

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return err, nil
	}

	adminRspBlk, err := buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return err, nil
	}

	return nil, adminRspBlk.CmdOutput
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest for loading a random value in one of the wrapping   */
/* key registers                                                              */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain where a random value is to be loaded  */
/*    in one of the wrapping key registers                                    */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* string -- the HTPRequest string with the signed CPRB for the command       */
/* error -- reports any errors                                                */
/*----------------------------------------------------------------------------*/
func CreateRandomWKReq(ci common.CommonInputs, de common.DomainEntry,
	sigkeys []string, sigkeySkis []string, sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_GEN_WK
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	// no input parameters
	return CreateSignedHTPRequest(ci, de, adminBlk, sigkeys, sigkeySkis, sigkeyTokens)
}
