//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 04/09/2021    CLH             Adapt for TKE SDK
// 07/30/2021    CLH             Add SSUrl to CommonInputs

package ep11cmds

import (
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Adds domain control points                                                 */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain whose attributes are to be set        */
/* []byte -- bit mask of control points to be enabled.  16 bytes are expected.*/
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func AddDomainControlPoints(ci common.CommonInputs, de common.DomainEntry,
	cpsToSet []byte, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) error {

	htpRequestString, err := AddDomainControlPointsReq(
		ci, de, cpsToSet, sigkeys, sigkeySkis, sigkeyTokens)
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
/* Creates the HTPRequest for adding domain control points                    */
/*----------------------------------------------------------------------------*/
func AddDomainControlPointsReq(ci common.CommonInputs, de common.DomainEntry,
	cpsToSet []byte, sigkeys []string, sigkeySkis []string, sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_DOM_CONTROLPOINT_ADD
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	adminBlk.CmdInput = cpsToSet
	return CreateSignedHTPRequest(ci, de, adminBlk, sigkeys, sigkeySkis, sigkeyTokens)
}
