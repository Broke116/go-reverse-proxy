package main

import (
	"net/http"
	"sync/atomic"
)

// HomeHandler returns information about proxy server
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reverse proxy server is up and running. Accepting at port " + port + " Redirecting to " + target1 + " , " + target2))
}

// HandleRequest is used for handling the incoming requests and reading the body of a request.
// then decodes the body content into requestPayload struct to extract proxy_condition value
// then it gets the proxyUrl depending on the proxy_condition value. finally it calls the serveReverseProxy function to redirect the request
func HandleRequest(w http.ResponseWriter, req *http.Request) {
	// url, err := getProxyURL(requestPayload.ProxyCondition)
	work := Work{id: atomic.AddUint64(&counter, 1), response: w, request: req}

	/*select {
	case requests <- work:
		logger.Printf("request with id %d will be redirected", work.id)
	case <-results:
		logger.Println("finished")
	}*/

	requests <- work

	select {
	case result := <-results:
		logger.Printf("request id %d was redirected", result.id)
	}
}
