package main

import (
	"container/heap"
	"fmt"
)

type Balancer struct {
	pool Pool
	done chan *Work
	i    int
}

// mapping redirection targets to each worker during the initialization
func setURL(wid int) string {
	urlMap := make(map[int]string)
	if wid%2 == 0 {
		urlMap[wid] = target2
	} else {
		urlMap[wid] = target1
	}
	return urlMap[wid]
}

func initBalancer() *Balancer {
	done := make(chan *Work, nWorkers)
	b := &Balancer{make(Pool, 0, nWorkers), done, 0}

	for i := 1; i <= nWorkers; i++ {
		url := setURL(i)
		w := &Work{wid: i, requests: make(chan *Request, capacity), url: url}
		heap.Push(&b.pool, w)
		go w.workIt(b.done)
	}
	return b
}

func (b *Balancer) balance(req chan *Request) {
	for {
		select {
		case request := <-req:
			b.dispatch(request)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req *Request) {
	w := heap.Pop(&b.pool).(*Work)
	w.requests <- req
	w.pending++
	fmt.Println()
	heap.Push(&b.pool, w)
	logger.Printf("started %d; pending now %d\n", w.wid, w.pending)
}

func (b *Balancer) completed(w *Work) {
	w.pending--
	heap.Remove(&b.pool, w.idx)
	heap.Push(&b.pool, w)
}
