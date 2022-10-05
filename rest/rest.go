package rest

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.uber.org/ratelimit"
)

type Rest interface {
	Get(uri string, options ...RequestOption) (*Response, error)
}

type void any

var signal void

type rest struct {
	requestOptions  []RequestOption
	client          *resty.Client
	concurrentLimit chan void
	rateLimit       ratelimit.Limiter
}

type Response struct {
	Body    []byte
	Headers http.Header
}

type RestOption func(*rest)

type RequestOption func(*resty.Request)

func (r *rest) Get(uri string, options ...RequestOption) (*Response, error) {
	req := r.client.R()
	if r.concurrentLimit != nil {
		r.concurrentLimit <- true
		defer func() {
			<-r.concurrentLimit
		}()
	}
	if r.rateLimit != nil {
		r.rateLimit.Take()
	}
	for _, o := range r.requestOptions {
		o(req)
	}
	for _, o := range options {
		o(req)
	}
	response, err := req.Get(uri)
	if err != nil {
		return nil, err
	}
	response.Time()
	sc := response.StatusCode()
	scOK := sc >= 200 && sc < 300
	if !scOK {
		if sc >= 500 {
			fmt.Println("Fatal Error")
		}
		return nil, fmt.Errorf("invalid status code %d, with %s", sc, response.Request.URL)
	}
	return &Response{
		Body:    response.Body(),
		Headers: response.Header(),
	}, nil
}

func New(options ...RestOption) Rest {
	r := &rest{
		client:         resty.New(),
		requestOptions: []RequestOption{},
	}
	for _, o := range options {
		o(r)
	}
	return r
}

// Rest Options

func WithDefaultRequestOptions(options ...RequestOption) RestOption {
	return func(r *rest) {
		r.requestOptions = options
	}
}

func WithBaseUrl(baseUrl string) RestOption {
	return func(r *rest) {
		r.client.SetBaseURL(baseUrl)
	}
}

func WithConcurrentLimit(size int) RestOption {
	return func(r *rest) {
		r.concurrentLimit = make(chan void, size)
	}
}

func WithRateLimit(rate int) RestOption {
	return func(r *rest) {
		rl := ratelimit.New(rate) // per second
		r.rateLimit = rl
	}
}

// RequestOption

func WithQueryParam(key, value string) RequestOption {
	return func(r *resty.Request) {
		r.SetQueryParam(key, value)
	}
}

func WithQueryParams(params map[string]string) RequestOption {
	return func(r *resty.Request) {
		r.SetQueryParams(params)
	}
}

func WithPathParam(key, value string) RequestOption {
	return func(r *resty.Request) {
		r.SetPathParam(key, value)
	}
}

func WithResult(result interface{}) RequestOption {
	return func(r *resty.Request) {
		r.SetResult(result)
	}
}
