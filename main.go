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

func main() {
	logger.Printf("Server will handle requests at %s\n", port)

	http.HandleFunc("/home", HomeHandler)

	http.HandleFunc("/", HandleRequest)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
