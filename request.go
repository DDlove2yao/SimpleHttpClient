package main

import (
	"bytes"
	"net/http"
)

type CustomRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Cookies []*http.Cookie
	Body    string
}

func (r *CustomRequest) Do() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBufferString(r.Body))
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
