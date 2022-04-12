package main

import (
	"errors"
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
	"strings"
	"time"
)

const (
	defaultContentType = "text/plain; charset=utf-8"
	HttpMethodPost     = "post"
	HttpMethodGet      = "get"
)

type WebHook struct {
	Url string

	httpClient *resty.Client
	username   string
	password   string
	method     string
	headers    map[string]string
}


type Plugin struct {
	Type     string            `yaml:"type"`
	Url      string            `yaml:"url"`
	Path     string            `yaml:"path"`
	Format   string            `yaml:"format"`
	User     string            `yaml:"user"`
	PassWord string            `yaml:"password"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
}

var config = &Plugin{
	Type:     "webhook",
	Url:      "https://sct.ftqq.com",
	Path:     "/",
	Format:   "",
	User:     "admin",
	PassWord: "admin",
	Method:   "get",
	Headers:  map[string]string{},
}

var items = map[string]string{"1": "hello ftp", "2": "kevin"}



func NewWebHookProvider(config *Plugin) *WebHook {
	url := fmt.Sprintf("%s%s", config.Url, config.Path)
	return &WebHook{
		Url:      url,
		username: config.User,
		password: config.PassWord,
		method:   strings.ToLower(config.Method),
		headers:  config.Headers,
	}
}

func (w *WebHook) Connect() error {
	client := resty.New()
	client.SetRESTMode()
	client.SetTimeout(20 * time.Second)
	client.SetContentLength(true)
	client.
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)

	if w.username != "" && w.password != "" {
		client.SetBasicAuth(w.username, w.password)
	}

	w.httpClient = client
	return nil
}

func (w *WebHook) Close() error {
	return nil
}

func (w *WebHook) Write(data map[string]string) error {
	var contentType string
	var err error
	for _, item := range data {
		var resp *resty.Response
		contentType = defaultContentType
		r := w.httpClient.R().SetBody(item).SetHeader("content-type", contentType)
		if len(w.headers) > 0 {
			r.SetHeaders(w.headers)
		}
		switch w.method {
		case HttpMethodPost:
			resp, err = r.Post(w.Url)
		case HttpMethodGet:
			resp, err = r.Get(w.Url)
		default:
			resp, err = r.Post(w.Url)
		}
		//fmt.Println(resp)
		if err != nil {
			return err
		}
		if (resp.StatusCode() != http.StatusOK) && (resp.StatusCode() != http.StatusCreated) {
			return errors.New(string(resp.Body()))
		} else {
			return nil
		}
	}
	return nil
}


func main() {
	f := NewWebHookProvider(config)
	_ = f.Connect()

	_ = f.Write(items)
	_ = f.Close()
}
