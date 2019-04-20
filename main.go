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
	address    = ":9090"
	urlA       = "http://localhost:4500"
	urlB       = "http://localhost:1332"
	defaultURL = "http://localhost:1333"
)

type requestPayload struct {
	ProxyCondition string `json:"proxy_condition"`
}

func initialLog() {
	logger.Printf("Redirecting to A url: %s\n", urlA)
	logger.Printf("Redirecting to B url: %s\n", urlB)
	logger.Printf("Redirecting to Default url: %s\n", defaultURL)
}

func logRequest(requestionPayload requestPayload, proxyURL string) {
	logger.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyURL)
}

func getProxyURL(proxyCondition string) (string, error) {
	var err error
	proxyCondition = strings.ToUpper(proxyCondition)

	if proxyCondition == "A" {
		if err = checkServerConnectivity(urlA); err == nil {
			return urlA, nil
		}
	}

	if proxyCondition == "B" {
		if err = checkServerConnectivity(urlB); err == nil {
			return urlB, nil
		}
	}

	if err = checkServerConnectivity(defaultURL); err == nil {
		return defaultURL, nil
	}

	return "", err
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
// then it gets the proxyUrl depending on the proxy_condition value. finally it calls the serveReverseProxy function to redirect the request
func handleRequests(res http.ResponseWriter, req *http.Request) {
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
	}

	logRequest(requestPayload, url)

	serveReserveProxy(url, res, req)
}

func main() {
	logger.Printf("Server will handle requests at %s\n", address)

	initialLog()

	http.HandleFunc("/", handleRequests)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
