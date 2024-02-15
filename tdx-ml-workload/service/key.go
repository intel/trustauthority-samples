/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	kbsclient "github.com/intel/kbs/v1/client"
	"github.com/intel/trustauthority-client/go-connector"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GetKeyRequest struct {
	AttestationToken string `json:"attestation_token"`
	KeyTransferUrl   string `json:"key_transfer_url"`
}

type GetKeyResponse struct {
	WrappedKey []byte `json:"wrapped_key"`
	WrappedSwk []byte `json:"wrapped_swk"`
}

func (t *GetKeyResponse) Headers() http.Header {
	return corsHeaders
}

func (mw loggingMiddleware) GetKey(ctx context.Context, req GetKeyRequest) (*GetKeyResponse, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("GetKey took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.GetKey(ctx, req)
	return resp, err
}

func (svc service) GetKey(_ context.Context, req GetKeyRequest) (*GetKeyResponse, error) {

	var err error
	var resp []byte
	var attestationType string
	var request *kbsclient.KeyTransferRequest

	keyUrl, _ := url.Parse(req.KeyTransferUrl)
	client := kbsclient.NewKBSClient(svc.httpClient, keyUrl, "")

	if req.AttestationToken != "" {
		request = &kbsclient.KeyTransferRequest{
			AttestationToken: req.AttestationToken,
		}
	} else {
		resp, attestationType, err = client.TransferKey()
		if err != nil {
			return nil, errors.Wrap(err, "could not get key")
		}

		var nonce *connector.VerifierNonce
		err = json.Unmarshal(resp, &nonce)
		if err != nil {
			return nil, errors.New("could not unmarshal ITA nonce")
		}

		quote, userData, err := collectEvidence(base64.StdEncoding.EncodeToString(append(nonce.Val, nonce.Iat[:]...)), svc.userData)
		if err != nil {
			return nil, errors.Wrap(err, "could not get quote")
		}

		request = &kbsclient.KeyTransferRequest{
			Quote:    quote,
			Nonce:    nonce,
			UserData: userData,
		}
	}

	resp, err = client.TransferKeyWithEvidence(request, attestationType)
	if err != nil {
		return nil, errors.Wrap(err, "could not transfer key")
	}

	var response GetKeyResponse
	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
