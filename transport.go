package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

//first we define our endpoint request and response types (plain go)
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

/*
 *	then we define our actual endpoints
 *	ussuallly if not always, endpoints wrap your service interface
 *	type Endpoint func(ctx, context, request interface{}) (response {interface}, err error)
 */

func makeUppecaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		req, ok := request.(uppercaseRequest)
		if !ok {
			return uppercaseResponse{"", "Bad Request"}, fmt.Errorf("Bad Request")

		}
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(countRequest)
		if !ok {
			return countResponse{0}, fmt.Errorf("Bad Request")
		}
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

/*
 * this functions are http transport decoders, take a http.Request and return a typed request for our service
 * they follow the type type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)
 */

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req countRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

/*
 * The following is a transport middleware
 * in gokit a middleware type is defined as: type Middleware func(Endpoint) Endpoint
 */

//this is a logging middleware (not in use)
func loginTransportMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "starting request")
			defer logger.Log("msg", "request finished")
			return next(ctx, request)
		}
	}
}
