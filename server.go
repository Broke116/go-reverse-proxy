package main

import (
	"log"
	"net"
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
		log.Printf("%s is unreachable %v", tcpAddress, err)
	}

	return err
}
