package ghttp

import (
	"net/http"
	"net/url"
	"strings"
)

func NewRequest(method, URL string, data, header map[string]string) (req *http.Request, err error) {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return NewPost(URL, data, header)
	default:
		return NewGet(URL, data, header)
	}
}

func NewGet(URL string, queryData, header map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", appendQuery(URL, queryData), nil)
	if err != nil {
		return
	}
	setHeader(req, header)
	return
}

func NewPost(URL string, data, header map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequest("POST", URL, strings.NewReader(encodeData(data)))
	if err != nil {
		return
	}
	setHeader(req, header)
	return
}

func appendQuery(URL string, queryData map[string]string) string {
	if len(queryData) == 0 {
		return URL
	}
	return URL + getConcatChar(URL) + encodeData(queryData)
}

func getConcatChar(URL string) string {
	if strings.Contains(URL, "?") {
		return "&"
	}
	return "?"
}

func encodeData(data map[string]string) string {
	values := url.Values{}
	if data != nil {
		for k, v := range data {
			values.Set(k, v)
		}
	}
	return values.Encode()
}

func setHeader(req *http.Request, header map[string]string) {
	if header == nil {
		return
	}
	foundContentType := false
	if header != nil {
		for k, v := range header {
			if strings.EqualFold(k, "host") {
				req.Host = v
				continue
			}
			req.Header.Set(k, v)
			if strings.TrimSpace(strings.ToLower(k)) == "content-type" {
				foundContentType = true
			}
		}
	}
	if !foundContentType {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}
