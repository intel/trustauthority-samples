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

func setProvisionHandler(svc service.Service, router *mux.Router, options []httpTransport.ServerOption) error {

	getProvisionHandler := httpTransport.NewServer(
		makeProvisionHTTPEndpoint(svc),
		decodeProvisionHTTPRequest,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/provision", getProvisionHandler).Methods(http.MethodPost)
	router.Handle("/provision", optionsHandler).Methods(http.MethodOptions)

	return nil
}

func makeProvisionHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.ProvisionRequest)
		return svc.Provision(ctx, req)
	}
}

func decodeProvisionHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {

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

	var req service.ProvisionRequest
	err := dec.Decode(&req)
	if err != nil {
		log.WithError(err).Error(ErrJsonDecodeFailed.Error())
		return nil, ErrJsonDecodeFailed
	}

	return req, nil
}
