/*
 * Copyright (c) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package client

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type KBSClient interface {
	TransferKey() ([]byte, string, error)
	TransferKeyWithEvidence(*KeyTransferRequest, string) ([]byte, error)
}

type kbsClient struct {
	Client  *http.Client
	BaseURL *url.URL
	ApiKey  string
}

func NewKBSClient(client *http.Client, baseURL *url.URL, apiKey string) KBSClient {
	return &kbsClient{
		Client:  client,
		BaseURL: baseURL,
		ApiKey:  apiKey,
	}
}

func (kc *kbsClient) requestAndProcessResponse(
	newRequest func() (*http.Request, error),
	queryParams map[string]string,
	headers map[string]string,
	processResponse func(*http.Response) error,
) error {
	var req *http.Request
	var err error

	if req, err = newRequest(); err != nil {
		return err
	}

	{
		if queryParams != nil {
			q := req.URL.Query()
			for param, val := range queryParams {
				q.Add(param, val)
			}
			req.URL.RawQuery = q.Encode()
		}
	}

	{
		for name, val := range headers {
			req.Header.Add(name, val)
		}
	}

	var resp *http.Response
	if resp, err = kc.Client.Do(req); err != nil {
		return err
	}

	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				logrus.WithError(err).Errorf("Failed to close response body")
			}
		}()
	}

	if resp.StatusCode != http.StatusOK || resp.ContentLength == 0 {
		return errors.Errorf("Invalid response: StatusCode = %d, ContentLength = %d", resp.StatusCode, resp.ContentLength)
	}

	return processResponse(resp)
}
