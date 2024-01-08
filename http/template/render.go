package gtemplate

import (
	"errors"
	ghash "github.com/snail007/gmc/util/hash"
	"sync"
)

var DefaultRender = NewRender()

type Render struct {
	tplCache map[string]*Template
	tplLock  sync.Mutex
	funcMap  map[string]interface{}
	leftDelims,
	rightDelims string
}

func NewRender() *Render {
	return &Render{
		tplCache: map[string]*Template{},
		funcMap:  map[string]interface{}{},
	}
}

func (s *Render) clearTemplateCache() {
	s.tplLock.Lock()
	defer s.tplLock.Unlock()
	s.tplCache = map[string]*Template{}
}

func (s *Render) Delims(left, right string) *Render {
	s.clearTemplateCache()
	s.leftDelims = left
	s.rightDelims = right
	return s
}

func (s *Render) AddFuncMap(m map[string]interface{}) *Render {
	if len(m) == 0 {
		return s
	}
	s.clearTemplateCache()
	for k, v := range m {
		s.funcMap[k] = v
	}
	return s
}

func (s *Render) getTemplate(tplBytes []byte) (*Template, error) {
	s.tplLock.Lock()
	defer s.tplLock.Unlock()
	id := ghash.Md5Bytes(tplBytes)
	if t, ok := s.tplCache[id]; ok {
		return t, nil
	}
	t := New()
	t.SetCtx(ctx)
	t.DdisableLogging()
	t.DisableLoadDefaultBinData()
	t.SetBinBytes(map[string][]byte{"tpl": tplBytes})
	if s.leftDelims != "" && s.rightDelims != "" {
		t.Delims(s.leftDelims, s.rightDelims)
	}
	if len(s.funcMap) > 0 {
		t.Funcs(s.funcMap)
	}
	err := t.Parse()
	if err != nil {
		return nil, err
	}
	s.tplCache[id] = t
	return t, nil
}

// Parse the template, tplBytesOrString is []byte or string template
func (s *Render) Parse(tplBytesOrString interface{}, tplData map[string]interface{}) (d []byte, err error) {
	var b []byte
	switch v := tplBytesOrString.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return nil, errors.New("wrong type of tpl data")
	}
	t, err := s.getTemplate(b)
	if err != nil {
		return
	}
	return t.Execute("tpl", tplData)
}
