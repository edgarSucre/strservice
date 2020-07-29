package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var svc StringService
	svc = stringService{}
	svc = logginMiddleware{logger, svc}

	/*
	 *	http handlers are part of the Transport layer, keeping it simple, just HTTP for now
	 * 	for this we use the package "github.com/go-kit/kit/transport/http"
	 */

	/*
		 * This section creates a transport middleware
		 * adding logging to the endpoint if needed
			uppercaseEndpoint := makeCountEndpoint(svc)
			loginEndpoint := loginTransportMiddleware(logger)(uppercaseEndpoint)
			uppercaseHandler := httptransport.NewServer(
				loginEndpoint,
				decodeUppercaseRequest,
				encodeResponse,
			)
	*/

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
	http.ListenAndServe(":8080", nil)
}
