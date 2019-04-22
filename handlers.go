package main

import (
	"net/http"
	"sync/atomic"
)

// Request is a unit of work which consists of a http request that are going to be handled by a worker and send it to the right target.
type Request struct {
	id       uint64
	response http.ResponseWriter
	request  *http.Request
}

// HomeHandler returns information about proxy server
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reverse proxy server is up and running. Accepting at port " + port + " Redirecting to " + target1 + " , " + target2))
}

// HandleRequest is used for handling the incoming requests and reading the body of a request.
// then decodes the body content into requestPayload struct to extract proxy_condition value
// then it gets the proxyUrl depending on the proxy_condition value. finally it calls the serveReverseProxy function to redirect the request
func HandleRequest(w http.ResponseWriter, req *http.Request) {
	request := Request{id: atomic.AddUint64(&counter, 1), response: w, request: req}

	requests <- request

	for {
		select {
		case result := <-results:
			logger.Println(result)
			return
		}
	}
}
