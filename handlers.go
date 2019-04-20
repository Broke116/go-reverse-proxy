package main

import "net/http"

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
	url := target1

	serveReserveProxy(url, w, req)
}
