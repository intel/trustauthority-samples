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

func setQuoteHandler(svc service.Service, router *mux.Router, options []httpTransport.ServerOption) error {

	getQuoteHandler := httpTransport.NewServer(
		makeGetQuoteHTTPEndpoint(svc),
		decodeGetQuoteHTTPRequest,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/quote", getQuoteHandler).Methods(http.MethodPost)

	return nil
}

func makeGetQuoteHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GetQuoteRequest)
		return svc.GetQuote(ctx, req)
	}
}

func decodeGetQuoteHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {

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

	var req service.GetQuoteRequest
	err := dec.Decode(&req)
	if err != nil {
		log.WithError(err).Error(ErrJsonDecodeFailed.Error())
		return nil, ErrJsonDecodeFailed
	}

	return req, nil
}
