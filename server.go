package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

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

func serveReserveProxy(target string, w http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url) // creating the reverse proxy

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host")) // identifying the originating IP address of a client
	req.Host = url.Host

	proxy.ServeHTTP(w, req)
}
