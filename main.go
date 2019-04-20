package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "main package ", log.LstdFlags|log.Lshortfile)

const (
	port    = ":9090"
	target1 = "http://localhost:1331"
	target2 = "http://localhost:1332"
)

type requestPayload struct {
	ProxyCondition string `json:"proxy_condition"`
}

func initialLog() {
	logger.Printf("Redirecting to A url: %s\n", target1)
	logger.Printf("Redirecting to B url: %s\n", target2)
}

func logRequest(requestionPayload requestPayload, proxyURL string) {
	logger.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyURL)
}

func getProxyURL(proxyCondition string) (string, error) {
	var err error
	proxyCondition = strings.ToUpper(proxyCondition)

	if proxyCondition == "A" {
		if err = checkServerConnectivity(target1); err == nil {
			return target1, nil
		}
	}

	if proxyCondition == "B" {
		if err = checkServerConnectivity(target2); err == nil {
			return target2, nil
		}
	}

	return "", err
}

func serveReserveProxy(target string, w http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url) // creating the reverse proxy

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host")) // identifying the originating IP address of a client
	req.Host = url.Host

	proxy.ServeHTTP(w, req)
}

// reading the body of a request. then decodes the body content into requestPayload struct to extract proxy_condition value
// then it gets the proxyUrl depending on the proxy_condition value. finally it calls the serveReverseProxy function to redirect the request
func handleRequests(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		logger.Printf("Error reading body: %v", err)
		panic(err)
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))

	var requestPayload requestPayload
	err = decoder.Decode(&requestPayload)

	if err != nil {
		logger.Printf("Request payload could not be decoded")
		panic(err)
	}
	url, err := getProxyURL(requestPayload.ProxyCondition)

	if err != nil {
		logger.Println("error", err)
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(err.Error()))
	} else {
		logRequest(requestPayload, url)

		serveReserveProxy(url, w, req)
	}
}

// HomeHandler returns information about proxy server
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reverse proxy server is up and running. Accepting at port " + port + " Redirecting to " + target1 + " , " + target2))
}

func main() {
	logger.Printf("Server will handle requests at %s\n", port)

	initialLog()

	http.HandleFunc("/home", HomeHandler)

	http.HandleFunc("/", handleRequests)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
