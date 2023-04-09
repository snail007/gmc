package ghttp

import (
	"github.com/snail007/gmc/util/gpool"
	gmap "github.com/snail007/gmc/util/map"
	"net/http"
	"sync"
)

type BatchRequest struct {
	reqArr            []*http.Request
	client            *http.Client
	waitFirstSuccess  bool
	pool              *gpool.GPool
	respArr           []*Response
	errArr            []error
	firstSuccessIndex int
	doFunc            func(idx int, req *http.Request) (*http.Response, error)
	afterExecute      func(*BatchRequest)
}

func NewBatchRequest(reqArr []*http.Request, client *http.Client) *BatchRequest {
	if client == nil {
		client = defaultClient()
	}
	return &BatchRequest{reqArr: reqArr, client: client}
}

func (s *BatchRequest) AfterExecute(afterExecute func(*BatchRequest)) *BatchRequest {
	s.afterExecute = afterExecute
	return s
}

func (s *BatchRequest) DoFunc(doFunc func(idx int, req *http.Request) (*http.Response, error)) *BatchRequest {
	s.doFunc = doFunc
	return s
}

// Success returns requests at least one is success in waitFirstSuccess mode.
//returns all requests are success in non waitFirstSuccess mode.
func (s *BatchRequest) Success() bool {
	if s.waitFirstSuccess {
		return s.respArr[0] != nil
	}
	for _, v := range s.errArr {
		if v != nil {
			return false
		}
	}
	return true
}

// Result returns the first response or first error when all request fail.
func (s *BatchRequest) Result() (*Response, error) {
	for _, v := range s.respArr {
		if v != nil {
			return v, nil
		}
	}
	for _, v := range s.errArr {
		if v != nil {
			return nil, v
		}
	}
	// this never occurred
	return nil, nil
}

// ResultAll returns all request's result,
// len(requests) = len(responses) = len(errros),
// responses[0] is the response of requests[0], it may be nil if requests[0] has an error.
// errros[0] is the error of requests[0], it may be nil if requests[0] has no error.
func (s *BatchRequest) ResultAll() ([]*Response, []error) {
	return s.respArr, s.errArr
}

// ErrorAll returns all request's non nil error,
func (s *BatchRequest) ErrorAll() (errors []error) {
	for _, v := range s.errArr {
		if v != nil {
			errors = append(errors, v)
		}
	}
	return
}

// ResponseAll returns all request's non nil response,
func (s *BatchRequest) ResponseAll() (responses []*Response) {
	for _, v := range s.respArr {
		if v != nil {
			responses = append(responses, v)
		}
	}
	return
}

// SetPool sets a *gpool.GPool to execute request. In default goroutine will be used.
func (s *BatchRequest) SetPool(pool *gpool.GPool) {
	s.pool = pool
}

// WaitFirstSuccess sets Execute() return immediately when get a success response.
func (s *BatchRequest) WaitFirstSuccess() *BatchRequest {
	s.waitFirstSuccess = true
	return s
}

func (s *BatchRequest) do(idx int, req *http.Request) (*http.Response, error) {
	if s.doFunc != nil {
		return s.doFunc(idx, req)
	}
	return s.client.Do(req)
}

// Execute batch send requests,
//	In default Execute will wait all request done. If you want to get the first success response,
// using BatchRequest.WaitFirstSuccess().Execute(), Execute() will return immediately when get a success response.
//	If all requests are fail, Execute() return after all requests done.
func (s *BatchRequest) Execute() *BatchRequest {
	defer func() {
		if s.afterExecute != nil {
			s.afterExecute(s)
		}
	}()
	if !s.waitFirstSuccess {
		// default wait all
		respMap := gmap.New()
		errMap := gmap.New()
		g := sync.WaitGroup{}
		g.Add(len(s.reqArr))
		for i := 0; i < len(s.reqArr); i++ {
			req := s.reqArr[i]
			idx := i
			worker := func() {
				defer g.Done()
				resp, err := s.do(idx, req)
				respMap.Store(idx, NewResponse(resp))
				errMap.Store(idx, err)
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
			respValue, _ := respMap.Load(i)
			if respValue == nil {
				s.respArr = append(s.respArr, nil)
			} else {
				s.respArr = append(s.respArr, respValue.(*Response))
			}
			errValue, _ := errMap.Load(i)
			if errValue == nil {
				s.errArr = append(s.errArr, nil)
			} else {
				s.errArr = append(s.errArr, errValue.(error))
			}
		}
	} else {
		errMap := gmap.New()
		doneChn := make(chan bool, 1)
		firstSuccessChn := make(chan *http.Response, 1)
		g := sync.WaitGroup{}
		g.Add(len(s.reqArr))
		for i := 0; i < len(s.reqArr); i++ {
			req := s.reqArr[i]
			idx := i
			worker := func() {
				defer g.Done()
				resp, err := s.do(idx, req)
				errMap.Store(idx, err)
				if err != nil {
					return
				}
				select {
				case firstSuccessChn <- resp:
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
				s.respArr = append(s.respArr, nil)
				v, _ := errMap.Load(i)
				if v == nil {
					s.errArr = append(s.errArr, nil)
				} else {
					s.errArr = append(s.errArr, v.(error))
				}
			}
		case resp := <-firstSuccessChn:
			// at least one request success
			s.respArr = append(s.respArr, NewResponse(resp))
			s.errArr = append(s.errArr, nil)
		}
	}
	return s
}
