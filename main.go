package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	sleepStr := os.Getenv("MILLISECOND")

	sleepMillisecond := 0
	if sleepStr != "" {
		sleepMillisecond, _= strconv.Atoi(sleepStr)
	}

	uri := os.Getenv("URI")
	if uri == "" {
		log.Fatalf("ENV URL must be set.")
	}
	go makeGetRequest(uri, sleepMillisecond)
	select {
	case <-time.After(60 * time.Second):
		log.Println("timeout 1")
	}
}

func makeGetRequest(uri string, sleepMillisecond int) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	client.Transport = transportWithTimeout(5 * time.Second)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("can't create request: %s", err.Error()), uri)
	}
	for {
		if sleepMillisecond > 0 {
			log.Printf("sleep %d Millisecond \n", sleepMillisecond)
			time.Sleep(time.Duration(sleepMillisecond)* time.Millisecond)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf(fmt.Sprintf("http get failed, err: %s, uri: %s", err.Error(), uri))
		}
		log.Println(fmt.Sprintf("http.status: %d", resp.StatusCode))
	}
}
func transportWithTimeout(connectTimeout time.Duration) http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   connectTimeout,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}


