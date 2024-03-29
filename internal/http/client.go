package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Scheme string

var (
	Http  Scheme = "http"
	Https Scheme = "https"
)

type Method string

func (m Method) String() string {
	return strings.ToUpper(string(m))
}

var (
	Get  Method = "get"
	Post Method = "post"
	Put  Method = "put"
)

type Request struct {
	Url    string
	Method Method
}

type Url struct {
	Scheme Scheme
	Domain string
	Path   string
}

func (u Url) String() string {
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Domain, u.Path)
}

type Response struct {
	Payload Payload
	Stats   Stats
}

type Stats struct {
	TotalTime time.Duration
}

type Payload struct {
	Stream io.Reader
	Data   string
}

type Client interface {
	Execute(Request) (Response, error)
}

type HttpClient struct {
	client *http.Client
}

func NewClient() Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	return HttpClient{
		client: httpClient,
	}
}

func (h HttpClient) Execute(req Request) (Response, error) {
	var err error
	var res *http.Response
	startTime := time.Now()
	switch req.Method {
	case Get:
		res, err = h.client.Get(req.Url)
	case Post:
		res, err = h.client.Post(req.Url, "application/json", bytes.NewBufferString(""))
	default:
		res, err = h.client.Get(req.Url)
	}
	if err != nil {
		return Response{}, err
	}
	body, _ := io.ReadAll(res.Body)
	return Response{
		Payload: Payload{
			Data: string(body),
		},
		Stats: Stats{
			TotalTime: time.Since(startTime),
		},
	}, nil
}
