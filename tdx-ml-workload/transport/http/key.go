/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/intel/trustauthority-samples/tdxexample/service"
	log "github.com/sirupsen/logrus"
)

func setKeyHandler(svc service.Service, router *mux.Router, options []httpTransport.ServerOption) error {

	getKeyHandler := httpTransport.NewServer(
		makeGetKeyHTTPEndpoint(svc),
		decodeGetKeyHTTPRequest,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/key", getKeyHandler).Methods(http.MethodPost)
	router.Handle("/key", optionsHandler).Methods(http.MethodOptions)

	return nil
}

func makeGetKeyHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GetKeyRequest)
		return svc.GetKey(ctx, req)
	}
}

func decodeGetKeyHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {

	if r.Header.Get(HTTPHeaderKeyAccept) != HTTPHeaderValueApplicationJson {
		log.Error(ErrInvalidAcceptHeader.Error())
		return nil, ErrInvalidAcceptHeader
	}

	if r.Header.Get(HTTPHeaderKeyContentType) != HTTPHeaderValueApplicationJson {
		log.Error(ErrInvalidContentTypeHeader.Error())
		return nil, ErrInvalidContentTypeHeader
	}

	if r.ContentLength == 0 {
		log.Error(ErrEmptyRequestBody.Error())
		return nil, ErrEmptyRequestBody
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req service.GetKeyRequest
	err := dec.Decode(&req)
	if err != nil {
		log.WithError(err).Error(ErrJsonDecodeFailed.Error())
		return nil, ErrJsonDecodeFailed
	}

	return req, nil
}
