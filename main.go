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
		return urlA
	}

	if proxyCondition == "B" {
		return urlB
	}

	return defaultURL
}

func serveReserveProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}

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
