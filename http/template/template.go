// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	"bytes"
	"encoding/base64"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	gotemplate "text/template"
)

var (
	defaultTpl = New()
)

//SetBinBase64 key is file path no slash prefix, value is file base64 encoded bytes contents.
func SetBinBase64(data map[string]string) {
	binData := map[string][]byte{}
	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic("init template bin data fail, error: " + err.Error())
		}
		binData[k] = b
	}
	defaultTpl.SetBinData(data)
}

//SetBinBytes key is file path no slash prefix, value is file's bytes contents.
func SetBinBytes(files map[string][]byte) {
	defaultTpl.SetBinBytes(files)
}

//SetBinString key is file path no slash prefix, value is file's string contents.
func SetBinString(files map[string]string) {
	defaultTpl.SetBinString(files)
}

type Template struct {
	rootDir string
	tpl     *gotemplate.Template
	parsed  bool
	ext     string
	ctx     gcore.Ctx
	binData map[string][]byte
}

func (s *Template) BinData() map[string][]byte {
	return s.binData
}

func (s *Template) SetBinBytes(binData map[string][]byte) {
	for k, v := range binData {
		s.binData[k] = v
	}
}

func (s *Template) SetBinString(binData map[string]string) {
	for k, v := range binData {
		s.binData[k] = []byte(v)
	}
}

func (s *Template) SetBinData(binData map[string]string) {
	for k, v := range binData {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic("init template bin data fail, error: " + err.Error())
		}
		s.binData[k] = b
	}
}

func New() (t *Template) {
	tpl := gotemplate.New("gmc").Option("missingkey=zero")
	t = &Template{
		tpl:     tpl,
		ext:     ".html",
		ctx:     gcore.ProviderCtx()(),
		binData: map[string][]byte{},
	}
	t.ctx.SetTemplate(t)
	return
}

// NewTemplate create a template object, and config it.
// rootDir is root path of view files folder.
func NewTemplate(ctx gcore.Ctx, rootDir string) (t *Template, err error) {
	absRootDir := ""
	if rootDir != "" {
		absRootDir, err = filepath.Abs(rootDir)
		if err != nil {
			return
		}
		absRootDir = strings.Replace(absRootDir, "\\", "/", -1)
	}
	tpl := gotemplate.New("gmc").Option("missingkey=zero")
	t = &Template{
		rootDir: absRootDir,
		tpl:     tpl,
		ext:     ".html",
		ctx:     ctx,
		binData: map[string][]byte{},
	}
	ctx.SetTemplate(t)
	return
}

// Delims sets the action delimiters to the specified strings, to be used in
// subsequent calls to Parse. Nested template
// definitions will inherit the settings. An empty delimiter stands for the
// corresponding default: {{ or }}.
// The return value is the template, so calls can be chained.
func (s *Template) Delims(left, right string) {
	s.tpl.Delims(left, right)
}

// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type or if the name cannot be used syntactically as a function in a template.
// It is legal to overwrite elements of the map. The return value is the template,
// so calls can be chained.
func (s *Template) Funcs(funcMap map[string]interface{}) {
	s.tpl.Funcs(funcMap)
}

func (s *Template) String() string {
	return s.tpl.DefinedTemplates()
}

//Extension sets template file extension, default is : .html
//only files have the extension will be parsed.
func (s *Template) Extension(ext string) {
	s.ext = ext
}

// Execute applies the view file associated with t that has the given name
// to the specified data object and return the output.
// If an error occurs executing the template,execution stops.
// A template may be executed safely in parallel.
func (s *Template) Execute(name string, data interface{}) (output []byte, err error) {
	name = strings.Replace(name, "\\", "/", -1)
	if strings.HasSuffix(name, s.ext) {
		name = strings.TrimSuffix(name, s.ext)
	}
	buf := &bytes.Buffer{}
	err = s.tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return
	}
	output = buf.Bytes()
	return
}

func (s *Template) clearBinData() {
	s.binData = map[string][]byte{}
}

// Parse load all view files data and parse it to internal template object.
// Mutiple call of Parse() only the first call worked.
func (s *Template) Parse() (err error) {
	if s.parsed {
		return
	}
	s.parsed = true
	s.Funcs(addFunc(s.ctx))
	if len(s.binData) > 0 {
		s.ctx.Logger().Infof("parse views from binary data")
		err = s.parseFromBinData()
	} else {
		s.ctx.Logger().Infof("parse views from disk")
		err = s.parseFromDisk()
	}
	if err != nil {
		return
	}
	s.clearBinData()
	return
}
func (s *Template) parseFromBinData() (err error) {
	for k, v := range s.binData {
		// template without extension
		html := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", k, string(v))
		_, err = s.tpl.Parse(html)
		if err != nil {
			return
		}
		// template with extension
		html = fmt.Sprintf("{{define \"%s\"}}%s{{end}}", k+s.ext, string(v))
		_, err = s.tpl.Parse(html)
		if err != nil {
			return
		}
	}
	return
}
func (s *Template) parseFromDisk() (err error) {
	names := []string{}
	err = s.tree(s.rootDir, &names)
	if err != nil {
		return
	}
	var b []byte
	for _, v := range names {
		b, err = ioutil.ReadFile(filepath.Join(s.rootDir, v+s.ext))
		if err != nil {
			return
		}

		// template without extension
		html := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", v, string(b))
		_, err = s.tpl.Parse(html)

		// template with extension
		html = fmt.Sprintf("{{define \"%s\"}}%s{{end}}", v+s.ext, string(b))
		_, err = s.tpl.Parse(html)
		if err != nil {
			return
		}
	}
	return
}
func (s *Template) tree(folder string, names *[]string) (err error) {
	f, err := os.Open(folder)
	if err != nil {
		return
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return
	}
	if !finfo.IsDir() {
		return
	}
	files, err := filepath.Glob(folder + "/*")
	if err != nil {
		return
	}
	var file *os.File
	for _, v := range files {
		if file != nil {
			file.Close()
		}
		file, err = os.Open(v)
		if err != nil {
			return
		}
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			err = s.tree(v, names)
			if err != nil {
				return err
			}
		} else {
			if !strings.HasSuffix(v, s.ext) {
				continue
			}
			v, err = filepath.Abs(v)
			if err != nil {
				return err
			}
			v0 := strings.Replace(v, "\\", "/", -1)
			v0 = strings.Replace(v0, s.rootDir+"/", "", -1)
			v0 = strings.TrimSuffix(v0, s.ext)
			*names = append(*names, v0)
		}
	}
	if file != nil {
		file.Close()
	}
	return
}

func (s *Template) Ctx() gcore.Ctx {
	return s.ctx
}

func (s *Template) SetCtx(ctx gcore.Ctx) {
	s.ctx = ctx
}

func (s *Template) Ext() string {
	return s.ext
}

func (s *Template) SetExt(ext string) {
	s.ext = ext
}

func (s *Template) Tpl() *gotemplate.Template {
	return s.tpl
}

func (s *Template) SetTpl(tpl *gotemplate.Template) {
	s.tpl = tpl
}

func (s *Template) RootDir() string {
	return s.rootDir
}

func (s *Template) SetRootDir(rootDir string) error {
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}
	absRootDir = strings.Replace(absRootDir, "\\", "/", -1)
	s.rootDir = absRootDir
	return nil
}
