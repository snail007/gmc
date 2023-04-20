package ghttp

import (
	"net/url"
	"strings"
)

type URLBuilder struct {
	*url.URL
	query   map[string]string
	holders []string
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
	if s.URL.Path == "" {
		s.URL.Path = "/"
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

func (s *URLBuilder) HostsURL(hosts []string) (urlArr []string) {
	h0 := s.URL.Host
	for _, h := range hosts {
		s.URL.Host = h
		urlArr = append(urlArr, s.String())
	}
	s.URL.Host = h0
	return
}

func (s *URLBuilder) PathsURL(paths []string) (urlArr []string) {
	p0 := s.URL.Path
	for _, p := range paths {
		s.URL.Path = p
		urlArr = append(urlArr, s.String())
	}
	s.URL.Path = p0
	return
}

func (s *URLBuilder) QueriesURL(queries []map[string]string) (urlArr []string) {
	q1 := s.query
	for _, q := range queries {
		s.Query(q)
		urlArr = append(urlArr, s.String())
	}
	s.query = q1
	return
}

func (s *URLBuilder) Holders(holders ...string) *URLBuilder {
	s.holders = append(s.holders, holders...)
	return s
}

func (s *URLBuilder) HolderValuesURL(holderValues ...[]string) (urlArr []string) {
	tpl := s.String()
	for i := range holderValues[0] {
		var oldNew []string
		for k := range s.holders {
			oldNew = append(oldNew, s.holders[k], holderValues[k][i])
		}
		urlArr = append(urlArr, strings.NewReplacer(oldNew...).Replace(tpl))
	}
	return
}
