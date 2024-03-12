/*
 * Copyright (C) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package http

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/intel/trustauthority-samples/tdxexample/service"
)

var optionsHandler = httpTransport.NewServer(
	makeOptionsHTTPEndpoint(),
	httpTransport.NopRequestDecoder,
	httpTransport.EncodeJSONResponse,
)

func makeOptionsHTTPEndpoint() endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return &service.OptionsResponse{}, nil
	}
}
