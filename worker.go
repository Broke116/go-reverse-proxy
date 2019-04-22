package main

// Work struct defines the each worker
type Work struct {
	wid     int
	idx     int // heap index
	work    chan Request
	pending int
	url     string
}

func (w *Work) workIt(done chan *Work) { /* id int, requests <-chan Request, results chan<- Result*/
	for {
		request := <-w.work
		logger.Printf("worker id %d ", w.wid)
		serveReserveProxy(done, w, request.id, request.response, request.request)
	}
}
