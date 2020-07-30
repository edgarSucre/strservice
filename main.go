package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	labels := []string{"method", "error"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "service_domain",
		Subsystem: "styring_service",
		Name:      "request_count",
		Help:      "Number of request received",
	}, labels)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "service_domain",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, labels)

	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "service_domain",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	var svc StringService
	svc = stringService{}
	//since our middleware implements StringService we can do:
	svc = logginMiddleware{logger, svc}

	//wiring middleware like we did for loggin
	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

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
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
