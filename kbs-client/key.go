/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/intel/trustauthority-client/go-connector"
)

const (
	HTTPHeaderKeyContentType       = "Content-Type"
	HTTPHeaderValueApplicationJson = "application/json"
	HTTPHeaderTypeXApiKey          = "x-api-key"
	HTTPHeaderKeyAccept            = "Accept"
	HTTPHeaderKeyAttestationType   = "Attestation-Type"
)

type KeyTransferRequest struct {
	AttestationToken string                   `json:"attestation_token,omitempty"`
	Quote            []byte                   `json:"quote,omitempty"`
	Nonce            *connector.VerifierNonce `json:"nonce,omitempty"`
	UserData         []byte                   `json:"user_data,omitempty"`
	EventLog         []byte                   `json:"event_log,omitempty"`
}

type KeyTransferResponse struct {
	WrappedKey string `json:"wrapped_key"`
	WrappedSwk string `json:"wrapped_swk"`
}

// TransferKey sends a POST request to Relying Party to retrieve the challenge data to be used as userdata for quote generation
func (kc *kbsClient) TransferKey() ([]byte, string, error) {

	newRequest := func() (*http.Request, error) {
		return http.NewRequest(http.MethodPost, kc.BaseURL.String(), nil)
	}

	var queryParams map[string]string = nil
	var headers = map[string]string{
		HTTPHeaderTypeXApiKey: kc.ApiKey,
		HTTPHeaderKeyAccept:   HTTPHeaderValueApplicationJson,
	}

	var body []byte
	var attestationType string
	processResponse := func(resp *http.Response) error {
		var err error
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		attestationType = resp.Header.Get(HTTPHeaderKeyAttestationType)
		return nil
	}

	if err := kc.requestAndProcessResponse(newRequest, queryParams, headers, processResponse); err != nil {
		return nil, "", err
	}

	return body, attestationType, nil
}

// TransferKeyWithEvidence sends a POST request to Relying Party to retrieve the actual resource
func (kc *kbsClient) TransferKeyWithEvidence(request *KeyTransferRequest, attestationType string) ([]byte, error) {

	newRequest := func() (*http.Request, error) {
		reqBytes, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}

		return http.NewRequest(http.MethodPost, kc.BaseURL.String(), bytes.NewReader(reqBytes))
	}

	var queryParams map[string]string = nil
	var headers = map[string]string{
		//	HTTPHeaderTypeXApiKey:        kc.ApiKey,
		HTTPHeaderKeyContentType: HTTPHeaderValueApplicationJson,
		HTTPHeaderKeyAccept:      HTTPHeaderValueApplicationJson,
		//	HTTPHeaderKeyAttestationType: attestationType,
	}

	var body []byte
	processResponse := func(resp *http.Response) error {
		var err error
		if body, err = io.ReadAll(resp.Body); err != nil {
			return err
		}
		return nil
	}

	if err := kc.requestAndProcessResponse(newRequest, queryParams, headers, processResponse); err != nil {
		return nil, err
	}

	return body, nil
}
