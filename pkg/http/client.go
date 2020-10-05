package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	http "net/http"
	"sync"
	"time"

	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/log"
	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"
	"github.com/sony/gobreaker"
)

type Client struct {
	cli     *gorequest.SuperAgent
	mu      sync.Mutex
	cb      *gobreaker.CircuitBreaker
	Headers http.Header
}

func NewClient() *Client {
	return &Client{
		cli: gorequest.New(),
	}
}

// SetTimeout sets a timeout for the requests
func (c *Client) SetTimeout(t time.Duration) {
	c.cli = c.cli.Timeout(t)
}

// UseCircuitBreaker turns on the circuit breaker on a given settings
func (c *Client) UseCircuitBreaker(st gobreaker.Settings) {
	c.cb = gobreaker.NewCircuitBreaker(st)
}

func (c *Client) send() ([]byte, error) {
	res, body, errs := c.cli.EndBytes()
	if errs != nil && len(errs) > 0 {
		return nil, errs[0]
	}

	msg := fmt.Sprintf("http client: response status %s", res.Status)

	if res.StatusCode != 200 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		log.Info(msg, &log.LogContext{
			"body": string(body),
		})
		return nil, errors.New(msg)
	}

	log.Info(msg, nil)
	return body, nil
}

// Get performs a HTTP GET request with a given struct as querystring and, if enabled,
// using the circuit breaker pattern
func (c *Client) Get(url string, querystring interface{}) ([]byte, error) {
	qs, err := query.Values(querystring)
	if err != nil {
		return nil, err
	}

	req := c.cli.Get(url).Query(qs.Encode())
	log.Info("http client", &log.LogContext{
		"method":    req.Method,
		"url":       req.Url,
		"queryData": req.QueryData,
	})

	c.mu.Lock()
	req.Header = c.Headers
	c.mu.Unlock()

	var body []byte

	// circuitbreaker is on
	if c.cb != nil {
		b, err := c.cb.Execute(func() (interface{}, error) {
			return c.send()
		})

		if err != nil {
			return nil, err
		}

		body, _ = b.([]byte)
		return body, nil
	}

	return c.send()
}
