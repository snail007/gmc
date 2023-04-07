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
	respArr           []*http.Response
	errArr            []error
	firstSuccessIndex int
}

func NewBatchRequest(reqArr []*http.Request, client *http.Client) *BatchRequest {
	if client == nil {
		client = defaultClient()
	}
	return &BatchRequest{reqArr: reqArr, client: client}
}

// Index returns the index of response that first success request in requests.
func (s *BatchRequest) Index() int {
	return s.firstSuccessIndex
}

// ErrorFirst returns the first request error.
func (s *BatchRequest) ErrorFirst() error {
	for _, v := range s.errArr {
		if v != nil {
			return v
		}
	}
	return nil
}

// Errors returns all request's error.
func (s *BatchRequest) Errors() (errors []error) {
	for _, v := range s.errArr {
		if v != nil {
			errors = append(errors, v)
		}
	}
	return
}

// ErrorLast returns the last request error.
func (s *BatchRequest) ErrorLast() error {
	if len(s.errArr) == 0 {
		return nil
	}
	for i := len(s.errArr) - 1; i >= 0; i++ {
		v := s.errArr[i]
		if v != nil {
			return v
		}
	}
	return nil
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

// Execute batch send requests, len(respArr) equals len(errArr),
// respArr[0] is the response of requests[0], it may be nil if requests[0] has an error.
// errArr[0] is the error of requests[0], it may be nil if requests[0] has no error.
//	In default Execute will wait all request done. If you want to get the first success response,
// using BatchRequest.WaitFirstSuccess().Execute(), Execute() will return immediately when get a success response.
//	If all requests are fail, Execute() return after all requests done.
func (s *BatchRequest) Execute() *BatchRequest {
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
				resp, err := s.client.Do(req)
				respMap.Store(idx, resp)
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
				s.respArr = append(s.respArr, respValue.(*http.Response))
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
				resp, err := s.client.Do(req)
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
			s.respArr = append(s.respArr, resp)
			s.errArr = append(s.errArr, nil)
		}
	}
	return s
}
