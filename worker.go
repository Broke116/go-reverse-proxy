package main

// Work struct defines the each worker
type Work struct {
	idx     int // heap index
	work    chan Request
	pending int
}

func worker(id int, requests <-chan Request, results chan<- Result) {
	for request := range requests {
		url := target1
		serveReserveProxy(request.id, url, request.response, request.request)
	}
}
