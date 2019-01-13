package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	address    = ":9090"
	urlA       = "http://localhost:1331"
	urlB       = "http://localhost:1332"
	defaultURL = "http://localhost:1333"
)

type requestPayload struct {
	ProxyCondition string `json:"proxy_condition"`
}

func initialLog() {
	log.Printf("Server will handle requests at %s\n", address)
	log.Printf("Redirecting to A url: %s\n", urlA)
	log.Printf("Redirecting to B url: %s\n", urlB)
	log.Printf("Redirecting to Default url: %s\n", defaultURL)
}

func logRequest(requestionPayload requestPayload, proxyURL string) {
	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyURL)
}

func getProxyURL(proxyCondition string) string {
	proxyCondition = strings.ToUpper(proxyCondition)

	if proxyCondition == "A" {
		if err := checkServerConnectivity(urlA); err == nil {
			return urlA
		}
	}

	if proxyCondition == "B" {
		if err := checkServerConnectivity(urlB); err == nil {
			return urlB
		}
	}

	return defaultURL
}

func serveReserveProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url) // creating the reverse proxy

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host")) // identifying the originating IP address of a client
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}

// reading the body of a request. then decodes the body content into requestPayload struct to extract proxy_condition value
// then it gets the proxyUrl depending on the proxy_condition value. finally it callls the serveReverseProxy function to redirect the request
func handleRequests(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))

	var requestPayload requestPayload
	err = decoder.Decode(&requestPayload)

	if err != nil {
		panic(err)
	}
	url := getProxyURL(requestPayload.ProxyCondition)
	logRequest(requestPayload, url)

	serveReserveProxy(url, res, req)
}

func main() {
	initialLog()

	http.HandleFunc("/", handleRequests)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
