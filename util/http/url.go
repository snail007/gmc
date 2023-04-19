package ghttp

import "net/url"

type URLBuilder struct {
	*url.URL
	query map[string]string
}

func NewURLBuilder() *URLBuilder {
	return &URLBuilder{
		URL:   &url.URL{},
		query: map[string]string{},
	}
}

func (s *URLBuilder) HTTP() *URLBuilder {
	s.URL.Scheme = "http"
	return s
}

func (s *URLBuilder) String() string {
	if s.URL.Scheme == "" {
		s.HTTP()
	}
	return AppendQuery(s.URL.String(), s.query)
}

func (s *URLBuilder) HTTPS() *URLBuilder {
	s.URL.Scheme = "https"
	return s
}

func (s *URLBuilder) Scheme(scheme string) *URLBuilder {
	s.URL.Scheme = scheme
	return s
}

func (s *URLBuilder) Host(host string) *URLBuilder {
	s.URL.Host = host
	return s
}

func (s *URLBuilder) Path(path string) *URLBuilder {
	s.URL.Path = path
	return s
}

func (s *URLBuilder) Query(data map[string]string) *URLBuilder {
	for k, v := range data {
		s.query[k] = v
	}
	return s
}
