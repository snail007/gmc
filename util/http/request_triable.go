package ghttp

import "net/http"

// TriableRequest when request fail, retry to request with MaxTryCount.
type TriableRequest struct {
	req    *http.Request
	client *http.Client
	maxTry int
}

// NewTriableRequest new a TriableRequest, maxTry is the max retry count when a request error occurred.
// If client is nil, default client will be used.
func NewTriableRequest(req *http.Request, client *http.Client, maxTry int) *TriableRequest {
	if client == nil {
		if client == nil {
			client = defaultClient()
		}
	}
	return &TriableRequest{req: req, client: client, maxTry: maxTry}
}

// Execute send request with retrying ability.
func (s *TriableRequest) Execute() (resp *http.Response, errs []error) {
	var err error
	for s.maxTry >= 0 {
		resp, err = s.client.Do(s.req)
		if err != nil {
			errs = append(errs, err)
			s.maxTry--
			continue
		}
		errs = nil
		break
	}
	return
}
