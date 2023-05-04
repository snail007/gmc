package ghttp

import (
	"bytes"
	"context"
	"github.com/snail007/gmc/util/gpool"
	gmap "github.com/snail007/gmc/util/map"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type BatchRequest struct {
	reqArr           []*http.Request
	reqBodyMap       *gmap.Map
	reqTimeoutMap    *gmap.Map
	client           *http.Client
	waitFirstSuccess bool
	pool             *gpool.GPool
	respArr          []*Response
	doFunc           func(idx int, req *http.Request) (*http.Response, error)
	beforeDo         []BeforeDoFunc
	afterDo          []AfterDoFunc
	maxTry           int
	checkErrorFunc   func(int, *http.Request, *http.Response) error
	once             sync.Once
	keepalive        bool
}

// NewBatchRequest new a BatchRequest by []*http.Request.
func NewBatchRequest(reqArr []*http.Request, client *http.Client) *BatchRequest {
	s := &BatchRequest{
		reqArr:        reqArr,
		client:        client,
		reqBodyMap:    gmap.New(),
		reqTimeoutMap: gmap.New(),
		keepalive:     true,
	}
	//init timeout
	for idx, r := range s.reqArr {
		if deadline, ok := r.Context().Deadline(); ok {
			s.reqTimeoutMap.Store(idx, deadline.Sub(time.Now()))
		}
	}
	return s
}

// NewBatchURL new a BatchRequest by url slice []string.
func NewBatchURL(client *http.Client, method string, urlArr []string, timeout time.Duration, data, header map[string]string) (tr *BatchRequest, err error) {
	var reqs []*http.Request
	var cancels []context.CancelFunc
	for _, v := range urlArr {
		r, cancel, e := NewRequest(method, v, timeout, data, header)
		if e != nil {
			return nil, e
		}
		cancels = append(cancels, cancel)
		reqs = append(reqs, r)
	}
	s := NewBatchRequest(reqs, client).
		AppendAfterDo(func(resp *Response) {
			cancels[resp.idx]()
		})
	return s, nil
}

func (s *BatchRequest) init() *BatchRequest {
	s.respArr = nil
	if s.client == nil && s.doFunc == nil {
		s.client = defaultClient()
	}
	s.once.Do(func() {
		//init body
		if s.maxTry > 0 {
			for idx, r := range s.reqArr {
				if IsFormRequest(r) && r.Body != nil {
					body, _ := ioutil.ReadAll(r.Body)
					s.reqBodyMap.Store(idx, body)
				}
				if deadline, ok := r.Context().Deadline(); ok {
					s.reqTimeoutMap.Store(idx, deadline.Sub(time.Now()))
				}
			}
		}
	})
	return s
}

// Keepalive sets enable or disable for request keepalive
func (s *BatchRequest) Keepalive(keepalive bool) *BatchRequest {
	s.keepalive = keepalive
	return s
}

// MaxTry sets the max retry count when a request error occurred.
func (s *BatchRequest) MaxTry(maxTry int) *BatchRequest {
	s.maxTry = maxTry
	return s
}

// CheckErrorFunc if returns an error, the request treat as fail.
func (s *BatchRequest) CheckErrorFunc(checkErrorFunc func(idx int, req *http.Request, resp *http.Response) error) *BatchRequest {
	s.checkErrorFunc = checkErrorFunc
	return s
}

// SetBeforeDo sets callback call before request sent.
func (s *BatchRequest) SetBeforeDo(beforeDo BeforeDoFunc) *BatchRequest {
	return s.setBeforeDo(beforeDo, true)
}

// AppendBeforeDo add a callback call before request sent.
func (s *BatchRequest) AppendBeforeDo(beforeDo BeforeDoFunc) *BatchRequest {
	return s.setBeforeDo(beforeDo, false)
}

func (s *BatchRequest) setBeforeDo(beforeDo BeforeDoFunc, isSet bool) *BatchRequest {
	if isSet {
		if beforeDo != nil {
			s.beforeDo = []BeforeDoFunc{beforeDo}
		} else {
			s.beforeDo = []BeforeDoFunc{}
		}
	} else if beforeDo != nil {
		s.beforeDo = append(s.beforeDo, beforeDo)
	}
	return s
}

// SetAfterDo sets callback call after request sent.
func (s *BatchRequest) SetAfterDo(afterDo AfterDoFunc) *BatchRequest {
	return s.setAfterDo(afterDo, true)
}

// AppendAfterDo add a callback call after request sent.
func (s *BatchRequest) AppendAfterDo(afterDo AfterDoFunc) *BatchRequest {
	return s.setAfterDo(afterDo, false)
}

func (s *BatchRequest) setAfterDo(afterDo AfterDoFunc, isSet bool) *BatchRequest {
	if isSet {
		if afterDo != nil {
			s.afterDo = []AfterDoFunc{afterDo}
		} else {
			s.afterDo = []AfterDoFunc{}
		}
	} else if afterDo != nil {
		s.afterDo = append(s.afterDo, afterDo)
	}
	return s
}

func (s *BatchRequest) callBeforeDo(idx int, req *http.Request) {
	for _, f := range s.beforeDo {
		f(idx, req)
	}
}

func (s *BatchRequest) callAfterDo(resp *Response) {
	for _, f := range s.afterDo {
		f(resp)
	}
}

// DoFunc sets a request sender.
func (s *BatchRequest) DoFunc(doFunc func(idx int, req *http.Request) (*http.Response, error)) *BatchRequest {
	s.doFunc = doFunc
	return s
}

// Success returns requests at least one is success in waitFirstSuccess mode.
// returns all requests are success in non waitFirstSuccess mode.
func (s *BatchRequest) Success() bool {
	if s.waitFirstSuccess {
		return s.respArr[0].Response != nil
	}
	for _, r := range s.respArr {
		if r.Err() != nil {
			return false
		}
	}
	return true
}

// Resp returns first success response
func (s *BatchRequest) Resp() *Response {
	for _, v := range s.respArr {
		if v.Response != nil {
			return v
		}
	}
	return nil
}

// RespAll returns all request's response, if you want to get if the response is success,
// checking if Response.Err() is nil.
func (s *BatchRequest) RespAll() (responseAll []*Response) {
	return s.respArr
}

// ErrorCount returns count of fail requests.
func (s *BatchRequest) ErrorCount() int {
	failCnt := 0
	for _, v := range s.respArr {
		if v.Err() != nil {
			failCnt++
		}
	}
	return failCnt
}

// Err returns first error of fail requests.
func (s *BatchRequest) Err() error {
	for _, v := range s.respArr {
		if v.Err() != nil {
			return v.Err()
		}
	}
	return nil
}

// ErrAll returns all errors of fail requests.
func (s *BatchRequest) ErrAll() (errs []error) {
	for _, v := range s.respArr {
		if v.Err() != nil {
			errs = append(errs, v.Err())
		}
	}
	return
}

// Pool sets a *gpool.GPool to execute request. In default goroutine will be used.
func (s *BatchRequest) Pool(pool *gpool.GPool) *BatchRequest {
	s.pool = pool
	return s
}

// WaitFirstSuccess sets Execute() return immediately when get a success response.
func (s *BatchRequest) WaitFirstSuccess() *BatchRequest {
	s.waitFirstSuccess = true
	return s
}

// Close close all response body can context cancel  func.
func (s *BatchRequest) Close() *BatchRequest {
	for _, v := range s.respArr {
		v.Close()
	}
	return s
}

func (s *BatchRequest) do(idx int, req *http.Request) (*http.Response, context.CancelFunc, error) {
	req.Close = !s.keepalive
	var setTimeout = func(req *http.Request) context.CancelFunc {
		if v, ok := s.reqTimeoutMap.Load(idx); ok {
			ctx, cancel := context.WithTimeout(context.Background(), v.(time.Duration))
			*req = *req.WithContext(ctx)
			return cancel
		}
		return nil
	}
	var call = func(idx int, req *http.Request) (*http.Response, context.CancelFunc, error) {
		var resp *http.Response
		var err error
		cancel := setTimeout(req)
		if s.doFunc != nil {
			resp, err = s.doFunc(idx, req)
		} else {
			resp, err = s.client.Do(req)
		}
		if cancel != nil && err != nil {
			cancel()
		}
		if err != nil {
			return nil, cancel, err
		}
		if s.checkErrorFunc != nil {
			err = s.checkErrorFunc(idx, req, resp)
			if err != nil {
				resp.Body.Close()
				return nil, cancel, err
			}
		}
		return resp, cancel, nil
	}
	if s.maxTry <= 0 {
		return call(idx, req)
	}
	maxTry := s.maxTry
	tryCount := 0
	var resp *http.Response
	var err error
	var cancel context.CancelFunc
	for tryCount <= maxTry {
		if v, ok := s.reqBodyMap.Load(idx); ok && v != nil && len(v.([]byte)) > 0 {
			req.Body = ioutil.NopCloser(bytes.NewReader(v.([]byte)))
		}
		resp, cancel, err = call(idx, req)
		if err != nil {
			tryCount++
			continue
		}
		break
	}
	return resp, cancel, err
}

// Execute batch send requests,
//
//	In default Execute will wait all request done. If you want to get the first success response,
//
// using BatchRequest.WaitFirstSuccess().Execute(), Execute() will return immediately when get a success response.
//
//	If all requests are fail, Execute() return after all requests done.
func (s *BatchRequest) Execute() *BatchRequest {
	s.init()
	if !s.waitFirstSuccess {
		// default wait all
		respMap := gmap.New()
		g := sync.WaitGroup{}
		g.Add(len(s.reqArr))
		for i := 0; i < len(s.reqArr); i++ {
			req := s.reqArr[i]
			idx := i
			worker := func() {
				defer g.Done()
				s.callBeforeDo(idx, req)
				startTime := time.Now()
				resp, cancel, err := s.do(idx, req)
				endTime := time.Now()
				if err != nil && cancel != nil {
					cancel()
				}
				response := &Response{
					idx:       idx,
					req:       req,
					Response:  resp,
					respErr:   err,
					usedTime:  endTime.Sub(startTime),
					startTime: startTime,
					endTime:   endTime,
					cancel:    cancel,
				}
				s.callAfterDo(response)
				respMap.Store(idx, response)
			}
			if s.pool != nil {
				s.pool.Submit(worker)
			} else {
				go worker()
			}
		}
		g.Wait()
		// all requests done, fill results.
		for i := 0; i < len(s.reqArr); i++ {
			resp, _ := respMap.Load(i)
			s.respArr = append(s.respArr, resp.(*Response))
		}
	} else {
		respMap := gmap.New()
		doneChn := make(chan bool, 1)
		firstSuccessChn := make(chan *Response, 1)
		g := sync.WaitGroup{}
		g.Add(len(s.reqArr))
		for i := 0; i < len(s.reqArr); i++ {
			req := s.reqArr[i]
			idx := i
			worker := func() {
				defer g.Done()
				s.callBeforeDo(idx, req)
				startTime := time.Now()
				resp, cancel, err := s.do(idx, req)
				endTime := time.Now()
				if err != nil && cancel != nil {
					cancel()
				}
				response := &Response{
					idx:       idx,
					req:       req,
					Response:  resp,
					respErr:   err,
					usedTime:  endTime.Sub(startTime),
					startTime: startTime,
					endTime:   endTime,
					cancel:    cancel,
				}
				s.callAfterDo(response)
				respMap.Store(idx, response)
				if err != nil {
					return
				}
				select {
				case firstSuccessChn <- response:
				default:
				}
			}
			if s.pool != nil {
				s.pool.Submit(worker)
			} else {
				go worker()
			}
		}
		go func() {
			g.Wait()
			select {
			case doneChn <- true:
			default:
			}
		}()
		select {
		case <-doneChn:
			//all requests fail
			for i := 0; i < len(s.reqArr); i++ {
				v, _ := respMap.Load(i)
				s.respArr = append(s.respArr, v.(*Response))
			}
		case resp := <-firstSuccessChn:
			// at least one request success
			s.respArr = append(s.respArr, resp)
		}
	}
	return s
}
