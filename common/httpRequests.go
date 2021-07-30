//
// Copyright 2021 IBM Inc. All rights reserved
// SPDX-License-Identifier: Apache2.0
//

// CHANGE HISTORY
//
// Date          Initials        Description
// 07/04/2021    CLH             Adapt for TKE SDK
// 07/23/2021    CLH             Fix URL for private endpoints
// 07/30/2021    CLH             Add SSUrl to CommonInputs

package common

import (
	"errors"
	"github.com/IBM/ibm-hpcs-tke-sdk/rest"
)

/*----------------------------------------------------------------------------*/
/* Determines the base URL to use for HTTP requests to the IBM Cloud          */
/*----------------------------------------------------------------------------*/
func getBaseURL(apiEndPoint string, region string) (string, error) {

	if apiEndPoint == "cloud.ibm.com" ||
		apiEndPoint == "https://cloud.ibm.com" {
		return "https://tke." + region + ".hs-crypto.cloud.ibm.com", nil
	} else if apiEndPoint == "test.cloud.ibm.com" ||
		apiEndPoint == "https://test.cloud.ibm.com" {
		return "https://tke." + region + ".hs-crypto.test.cloud.ibm.com", nil
	} else if apiEndPoint == "private.cloud.ibm.com" ||
		apiEndPoint == "https://private.cloud.ibm.com" {
		return "https://tke.private." + region + ".hs-crypto.cloud.ibm.com", nil
	} else if apiEndPoint == "private.test.cloud.ibm.com" ||
		apiEndPoint == "https://private.test.cloud.ibm.com" {
		return "https://tke.private." + region + ".hs-crypto.test.cloud.ibm.com", nil
	} else {
		return "", errors.New("Invalid API endpoint")
	}
}

/*----------------------------------------------------------------------------*/
/* Creates the HTTP request for querying the domains for a crypto instance.   */
/*----------------------------------------------------------------------------*/
func CreateGetHsmsRequest(ci CommonInputs) (*rest.Request, error) {

	urlStart, err := getBaseURL(ci.ApiEndpoint, ci.Region)
	if err != nil {
		return nil, err
	}

	url := urlStart + "/v1/tke/" + ci.InstanceId + "/hsms"
	req := rest.GetRequest(url)
	req.Set("Content-type", "application/json")
	req.Set("Authorization", ci.AuthToken)
	return req, nil
}

/*----------------------------------------------------------------------------*/
/* Creates the HTTP request for sending an HTPRequest to a TKE catcher        */
/* program.                                                                   */
/*----------------------------------------------------------------------------*/
func CreatePostHsmsRequest(ci CommonInputs, hsmId string, htpRequest string) (*rest.Request, error) {

	urlStart, err := getBaseURL(ci.ApiEndpoint, ci.Region)
	if err != nil {
		return nil, err
	}

	url := urlStart + "/v1/tke/" + ci.InstanceId + "/hsms/" + hsmId
	req := rest.PostRequest(url)
	req.Set("Content-type", "application/json")
	req.Set("Authorization", ci.AuthToken)
	req.Body(`{"request":"` + htpRequest + `"}`)
	return req, nil
}

/*----------------------------------------------------------------------------*/
/* Creates an HTTP request to a signing service specified by the user to      */
/* return the public part of a signature key                                  */
/*----------------------------------------------------------------------------*/
func CreateGetPublicKeyRequest(sigkeyToken string, ssURL string,
		sigkey string) *rest.Request {

	url := ssURL + "/keys/" + sigkey
	req := rest.GetRequest(url)
	req.Set("Content-type", "application/json")
	req.Set("Authorization", sigkeyToken)
	return req
}

/*----------------------------------------------------------------------------*/
/* Creates an HTTP request to a signing service specified by the user to      */
/* sign data using a signature key                                            */
/*----------------------------------------------------------------------------*/
func CreateSignDataRequest(sigkeyToken string, ssURL string,
		sigkey string, dataToSign string) *rest.Request {

	url := ssURL + "/sign/" + sigkey
	req := rest.PostRequest(url)
	req.Set("Content-type", "application/json")
	req.Set("Authorization", sigkeyToken)
	req.Body(`{"hash_algorithm":"sha2-512","input":"` + dataToSign + `"}`)
	return req
}
