package base

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

func init() {
	//http client
	var (
		connectTimeout = 30 * time.Second
		headerTimeout  = 60 * time.Second
		keepAlive      = 60 * time.Second
	)
	httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				//connectTimeout
				Timeout: connectTimeout,
				//keepAlive
				KeepAlive: keepAlive,
			}).DialContext,
			MaxIdleConnsPerHost: 200,
			//keepAlive
			IdleConnTimeout: keepAlive,
			//headerTimeout
			ResponseHeaderTimeout: headerTimeout,
		},
	}
}

//CURL2Reader CURL基础
func CURL2Reader(addr, method string, headers map[string]string, body io.Reader, buffer *bytes.Buffer, exitChan <-chan bool) (map[string]interface{}, error) {
	//readTimeout
	var readTimeout = 300 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	go func() {
		for {
			exit, ok := <-exitChan
			if !ok {
				break
			}
			if exit {
				cancel()
				break
			}
		}
	}()
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		if req.Header.Get(k) != "" {
			req.Header.Set(k, v)
		} else {
			req.Header.Add(k, v)
		}
	}
	resp, err := httpClient.Do(req.WithContext(ctx))
	if resp != nil {
		defer func() {
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()
	}
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"StatusCode": resp.StatusCode,
	}
	for k, v := range resp.Header {
		result[k] = v[0]
	}
	return result, nil
}
