package main

import (
	"log"
	"net/http"
	"os"
)

const (
	port    = ":9090"
	target1 = "http://localhost:1331"
	target2 = "http://localhost:1332"
)

var logger = log.New(os.Stdout, "main package ", log.LstdFlags|log.Lshortfile)
var counter uint64 // request counter for now. in the future it will be changed
var requests = make(chan Request, 100)
var results = make(chan Result, 100)

func main() {
	logger.Printf("Server will handle requests at %s\n", port)

	for w := 1; w <= 4; w++ {
		go worker(w, requests, results)
	}

	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/", HandleRequest)

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
