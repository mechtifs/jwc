package requests

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type SessionOptions struct {
	Header     http.Header
	DoRedirect bool
	Timeout    time.Duration
}

type Session struct {
	http.Client
	Header http.Header
}

func NewSession(opt *SessionOptions) *Session {
	s := &Session{}
	if opt != nil {
		if opt.Header != nil {
			s.Header = opt.Header
		}
		if !opt.DoRedirect {
			s.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if opt.Timeout != 0 {
			s.Timeout = opt.Timeout
		}
	}
	return s
}

func (s *Session) Request(method string, url string, opt *RequestOptions) (*Response, error) {
	req, err := http.NewRequest(method, url, nil)
	req.Header = s.Header
	if err != nil {
		return nil, err
	}

	if opt != nil {
		if opt.Header != nil {
			for k, v := range opt.Header {
				req.Header[k] = v
			}
		}
		if opt.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(opt.Body))
		}
	}

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &Response{
		Url:        url,
		Content:    buf,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}, nil
}

func (s *Session) Get(url string, opt *RequestOptions) (*Response, error) {
	return s.Request("GET", url, opt)
}

func (s *Session) Post(url string, opt *RequestOptions) (*Response, error) {
	return s.Request("POST", url, opt)
}

func (s *Session) Put(url string, opt *RequestOptions) (*Response, error) {
	return s.Request("PUT", url, opt)
}

func (s *Session) Delete(url string, opt *RequestOptions) (*Response, error) {
	return s.Request("DELETE", url, opt)
}

func (s *Session) Trace(url string, opt *RequestOptions) (*Response, error) {
	return s.Request("TRACE", url, opt)
}

func (s *Session) Connect(url string, opt *RequestOptions) (*Response, error) {
	return s.Request("CONNECT", url, opt)
}
