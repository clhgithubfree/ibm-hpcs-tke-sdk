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
	"encoding/binary"

	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Reads an OA certificate                                                    */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the crypto module and domain to be queried       */
/* certificateIndex -- index into the certificate chain                       */
/*    0 = currently active epoch key, 1 = its parent, etc.                    */
/*                                                                            */
/* Outputs:                                                                   */
/* []byte -- the returned OA certificate                                      */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryDeviceCertificate(ci common.CommonInputs, de common.DomainEntry,
	certificateIndex uint32) ([]byte, error) {

	htpRequestString := QueryDeviceCertificateReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex(), certificateIndex)

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
/* Creates the HTPRequest to return a specific OA certificate                 */
/*                                                                            */
/* Inputs:                                                                    */
/* cryptoModuleIndex -- identifies the crypto module to be queried            */
/* domainIndex -- the domain assigned to the user.  We know this is a control */
/*    domain, and it will get past the domain check when the cloud processes  */
/*    the POST /hsms request.                                                 */
/* certificateIndex -- index into the certificate chain                       */
/*    0 = currently active epoch key, 1 = its parent, etc.                    */
/*                                                                            */
/* Output:                                                                    */
/* string -- hexadecimal string representing the HTPRequest                   */
/*----------------------------------------------------------------------------*/
func QueryDeviceCertificateReq(cryptoModuleIndex int, domainIndex int,
	certificateIndex uint32) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DEVICE_CERT
	adminBlk.DomainID = XCP_DOMAIN_0
	// module ID not used for queries
	// transaction counter not used for queries
	adminBlk.CmdInput = common.Uint32To4ByteSlice(certificateIndex)
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}

/*----------------------------------------------------------------------------*/
/* Returns the number of OA certificates in the OA certificate chain          */
/*                                                                            */
/* Inputs:                                                                    */
/* CommonInputs -- A structure containing inputs needed for all TKE SDK       */
/*      functions.  This includes: the API endpoint and region, the HPCS      */
/*      service instance id, an IBM Cloud authentication token, and the       */
/*      URL and port for the signing service if one is used.                  */
/* DomainEntry -- identifies the crypto module and domain to be queried       */
/*                                                                            */
/* Outputs:                                                                   */
/* uint32 -- the number of certificates in the OA certificate chain           */
/* error -- reports any errors for the operation                              */
/*----------------------------------------------------------------------------*/
func QueryNumberDeviceCertificates(ci common.CommonInputs,
	de common.DomainEntry) (uint32, error) {

	htpRequestString := QueryNumberDeviceCertificatesReq(
		de.GetCryptoModuleIndex(), de.GetDomainIndex())

	req, err := common.CreatePostHsmsRequest(ci, de.Hsm_id, htpRequestString)
	if err != nil {
		return 0, err
	}

	htpResponseString, err := common.SubmitHTPRequest(req)
	if err != nil {
		return 0, err
	}

	adminRspBlk, err := buildAdminRspBlk(htpResponseString, de)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(adminRspBlk.CmdOutput), nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTPRequest to return the number of OA certificates in the OA   */
/* certificate chain                                                          */
/*                                                                            */
/* Inputs:                                                                    */
/* cryptoModuleIndex -- identifies the crypto module to be queried            */
/* domainIndex -- the domain assigned to the user.  We know this is a control */
/*    domain, and it will get past the domain check when the cloud processes  */
/*    the POST /hsms request.                                                 */
/*                                                                            */
/* Output:                                                                    */
/* string -- hexadecimal string representing the HTPRequest                   */
/*----------------------------------------------------------------------------*/
func QueryNumberDeviceCertificatesReq(cryptoModuleIndex int,
	domainIndex int) string {

	var adminBlk AdminBlk
	adminBlk.CmdID = XCP_ADMQ_DEVICE_CERT
	adminBlk.DomainID = XCP_DOMAIN_0
	// module ID not used for queries
	// transaction counter not used for queries
	// empty payload to get the number of certificates
	return CreateQueryHTPRequest(cryptoModuleIndex, domainIndex, adminBlk)
}
