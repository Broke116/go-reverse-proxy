package main

import "net/http"

var requests = make(chan Work, 100)
var results = make(chan Result, 100)

// Work is a unit of work which consists of a http request that are going to be handled by a worker and send it to the right target.
type Work struct {
	id       uint64
	response http.ResponseWriter
	request  *http.Request
}

// Result is a struct which indicates the result of the redirection.
type Result struct {
	id     uint64
	result bool
}

func worker(id int, requests <-chan Work, results chan<- Result) {
	for work := range requests {
		url := target1
		serveReserveProxy(work.id, url, work.response, work.request)

		//logger.Printf("request id %d was redirected to  %s\n", work.id, url)
	}
}
