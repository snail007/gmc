package gurl

import (
	"net/url"
	"strings"
)

type Builder struct {
	*url.URL
	query   map[string]string
	holders []string
}

func NewBuilder() *Builder {
	return &Builder{
		URL:   &url.URL{},
		query: map[string]string{},
	}
}

func (s *Builder) HTTP() *Builder {
	s.URL.Scheme = "http"
	return s
}

func (s *Builder) String() string {
	if s.URL.Scheme == "" {
		s.HTTP()
	}
	if s.URL.Path == "" {
		s.URL.Path = "/"
	}
	return AppendQuery(s.URL.String(), s.query)
}

func (s *Builder) HTTPS() *Builder {
	s.URL.Scheme = "https"
	return s
}

func (s *Builder) Scheme(scheme string) *Builder {
	s.URL.Scheme = scheme
	return s
}

func (s *Builder) Host(host string) *Builder {
	s.URL.Host = host
	return s
}

func (s *Builder) Path(path string) *Builder {
	s.URL.Path = path
	return s
}

func (s *Builder) Query(data map[string]string) *Builder {
	for k, v := range data {
		s.query[k] = v
	}
	return s
}

func (s *Builder) HostsURL(hosts []string) (urlArr []string) {
	h0 := s.URL.Host
	for _, h := range hosts {
		s.URL.Host = h
		urlArr = append(urlArr, s.String())
	}
	s.URL.Host = h0
	return
}

func (s *Builder) PathsURL(paths []string) (urlArr []string) {
	p0 := s.URL.Path
	for _, p := range paths {
		s.URL.Path = p
		urlArr = append(urlArr, s.String())
	}
	s.URL.Path = p0
	return
}

func (s *Builder) QueriesURL(queries []map[string]string) (urlArr []string) {
	q1 := s.query
	for _, q := range queries {
		s.Query(q)
		urlArr = append(urlArr, s.String())
	}
	s.query = q1
	return
}

func (s *Builder) Holders(holders ...string) *Builder {
	s.holders = append(s.holders, holders...)
	return s
}

func (s *Builder) HolderValuesURL(holderValues ...[]string) (urlArr []string) {
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

func GetConcatChar(URL string) string {
	if strings.Contains(URL, "?") {
		return "&"
	}
	return "?"
}

func AppendQuery(URL string, queryData map[string]string) string {
	if len(queryData) == 0 {
		return URL
	}
	return URL + GetConcatChar(URL) + EncodeData(queryData)
}

func EncodeData(data map[string]string) string {
	values := url.Values{}
	if data != nil {
		for k, v := range data {
			values.Set(k, v)
		}
	}
	return values.Encode()
}
