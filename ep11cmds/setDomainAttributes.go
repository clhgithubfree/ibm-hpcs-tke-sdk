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
/* Sets the domain attributes                                                 */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the domain whose attributes are to be set        */
/* DomainAttributes -- new set of attributes to be loaded in the domain       */
/* []string -- identifies the signature keys to use to sign the command       */
/* []string -- the Subject Key Identifiers for the signature keys             */
/* []string -- authentication tokens for the signature keys                   */
/*                                                                            */
/* Outputs:                                                                   */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func SetDomainAttributes(ci common.CommonInputs,
	de common.DomainEntry, newAttributes DomainAttributes,
	sigkeys []string, sigkeySkis []string, sigkeyTokens []string) error {

	htpRequestString, err := SetDomainAttributesReq(
		ci, de, newAttributes, sigkeys, sigkeySkis, sigkeyTokens)
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
/* Creates the HTPRequest for setting the domain attributes                   */
/*----------------------------------------------------------------------------*/
func SetDomainAttributesReq(ci common.CommonInputs, de common.DomainEntry,
	newAttributes DomainAttributes, sigkeys []string, sigkeySkis []string,
	sigkeyTokens []string) (string, error) {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADM_DOM_SET_ATTR
	// administrative domain filled in later
	// module ID filled in later
	// transaction counter filled in later
	// assemble the payload
	adminBlk.CmdInput = make([]byte, 4*8)
	copy(adminBlk.CmdInput[0:4], []byte{0x00, 0x00, 0x00, 0x01})
	copy(adminBlk.CmdInput[4:8], common.Uint32To4ByteSlice(newAttributes.SignatureThreshold))
	copy(adminBlk.CmdInput[8:12], []byte{0x00, 0x00, 0x00, 0x02})
	copy(adminBlk.CmdInput[12:16], common.Uint32To4ByteSlice(newAttributes.RevocationSignatureThreshold))
	copy(adminBlk.CmdInput[16:20], []byte{0x00, 0x00, 0x00, 0x03})
	copy(adminBlk.CmdInput[20:24], common.Uint32To4ByteSlice(newAttributes.Permissions))
	copy(adminBlk.CmdInput[24:28], []byte{0x00, 0x00, 0x00, 0x04})
	copy(adminBlk.CmdInput[28:32], common.Uint32To4ByteSlice(newAttributes.OperationalMode))

	return CreateSignedHTPRequest(ci, de, adminBlk, sigkeys, sigkeySkis, sigkeyTokens)
}
