package ghttp

import (
	"github.com/snail007/gmc/util/gpool"
	gmap "github.com/snail007/gmc/util/map"
	"net/http"
	"sync"
	"time"
)

type BatchRequest struct {
	reqArr            []*http.Request
	client            *http.Client
	waitFirstSuccess  bool
	pool              *gpool.GPool
	respArr           []*Response
	firstSuccessIndex int
	doFunc            func(idx int, req *http.Request) (*http.Response, error)
	afterExecute      func(*BatchRequest)
}

func NewBatchRequest(reqArr []*http.Request, client *http.Client) *BatchRequest {
	return &BatchRequest{reqArr: reqArr, client: client}
}

func (s *BatchRequest) init() *BatchRequest {
	if s.client == nil {
		s.client = defaultClient()
	}
	return s
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

// RespAll returns all request's response,
func (s *BatchRequest) RespAll() (responses []*Response) {
	return s.respArr
}

func (s *BatchRequest) ErrorCount() int {
	failCnt := 0
	for _, v := range s.respArr {
		if v.Err() != nil {
			failCnt++
		}
	}
	return failCnt
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
	s.init()
	defer func() {
		if s.afterExecute != nil {
			s.afterExecute(s)
		}
	}()
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
				startTime := time.Now()
				resp, err := s.do(idx, req)
				endTime := time.Now()
				response := &Response{
					idx:       idx,
					req:       req,
					Response:  resp,
					respErr:   err,
					usedTime:  endTime.Sub(startTime),
					startTime: startTime,
					endTime:   endTime,
				}
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
				startTime := time.Now()
				resp, err := s.do(idx, req)
				endTime := time.Now()
				response := &Response{
					idx:       idx,
					req:       req,
					Response:  resp,
					respErr:   err,
					usedTime:  endTime.Sub(startTime),
					startTime: startTime,
					endTime:   endTime,
				}
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
