package main

import (
	"fmt"
	"gopkg.in/resty.v1"
	"time"
)

var Headers = map[string]string{"1": "ck_god", "2": "god_girl"}

var Body = []byte("hello,I am a http client!")

const (
	URL = "http://127.0.0.1:8082/rush/v1/hello/3?q=1"
)

func main() {
	var httpClient *resty.Client
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(6 * time.Second))
	client.SetContentLength(true)
	// Headers for all request
	client.SetHeaders(Headers)
	client.
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(5 * time.Second)).
		SetRetryMaxWaitTime(20 * time.Second)

	httpClient = client

	resp, err := httpClient.R().
		SetBody(Body).
		Get(URL)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
