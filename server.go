package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func checkServerConnectivity(address string) error {
	url, _ := url.Parse(address)
	urlPort := url.Port()
	urlHostname := url.Hostname()

	if urlHostname == "localhost" {
		urlHostname = "127.0.0.1"
	}

	tcpAddress := urlHostname + ":" + urlPort

	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", tcpAddress, timeout)

	if err != nil {
		logger.Printf("%s is unreachable %v", tcpAddress, err)
	}

	return err
}

func serveReserveProxy(requestID uint64, target string, w http.ResponseWriter, req *http.Request) {
	err := checkServerConnectivity(target)

	if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(err.Error()))
		result := Result{id: requestID, result: false}
		results <- result
	}

	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url) // creating the reverse proxy

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host")) // identifying the originating IP address of a client
	req.Host = url.Host

	proxy.ServeHTTP(w, req)

	result := Result{id: requestID, result: true}
	results <- result
}
