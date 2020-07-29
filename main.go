package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

/*
*	Service Layer (plain go)
 */

//StringService interface provides an  apbstraction of the business logic
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

//StringService implementation
type stringService struct{}

//ErrEmpty is returned when the input string is empty
var ErrEmpty = errors.New("Empty string")

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

/*
*	Endpoint Layer
 */

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
*	Transport layer, keeping it simple, just HTTP for now
* 	for this we use the package "github.com/go-kit/kit/transport/http"
 */

func main() {
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		makeUppecaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
