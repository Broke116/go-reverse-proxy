package main

import "container/heap"

type Balancer struct {
	pool Pool
	done chan *Work
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
	b := &Balancer{make(Pool, 0, nWorkers), done}

	for i := 0; i < nWorkers; i++ {
		url := setURL(i)
		w := &Work{wid: i, work: make(chan Request, capacity), url: url}
		heap.Push(&b.pool, w)
		go w.workIt(b.done)
	}
	return b
}

func (b *Balancer) balance(req chan Request) {
	for {
		select {
		case request := <-req:
			b.dispatch(request)
		case w := <-b.done:
			//results <- true
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	w := heap.Pop(&b.pool).(*Work)
	w.work <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Work) {
	w.pending--
	heap.Remove(&b.pool, w.idx)
	heap.Push(&b.pool, w)
}
