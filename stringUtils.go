package main

import (
	"io"
	"net/http"
	"regexp"
)

func CombineUrl(protocol, url string) string {
	if regexp.MustCompile(`^(localhost|\d+\.\d+\.\d+\.\d+|([a-zA-Z0-9\-\.]+))(:\d+)?$`).MatchString(url) {
		return protocol + "://" + url
	}
	return ""
}

func readResponseBody(resp *http.Response) string {
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(body)
}
