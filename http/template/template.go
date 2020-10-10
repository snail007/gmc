// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmctemplate

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	gotemplate "text/template"
)

var (
	bindata = map[string][]byte{}
)

func SetBinData(data map[string]string) {
	bindata = map[string][]byte{}
	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic("init template bin data fail, error: " + err.Error())
		}
		bindata[k] = b
	}
}

type Template struct {
	rootDir string
	tpl     *gotemplate.Template
	parsed  bool
	ext     string
}

// New create a template object, and config it.
// rootDir is root path of view files folder.
func New(rootDir string) (t *Template, err error) {
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return
	}
	absRootDir = strings.Replace(absRootDir, "\\", "/", -1)
	tpl := gotemplate.New("gmc").Option("missingkey=zero")
	t = &Template{
		rootDir: absRootDir,
		tpl:     tpl,
		ext:     ".html",
	}
	t.Funcs(addFunc())
	return
}

// Delims sets the action delimiters to the specified strings, to be used in
// subsequent calls to Parse. Nested template
// definitions will inherit the settings. An empty delimiter stands for the
// corresponding default: {{ or }}.
// The return value is the template, so calls can be chained.
func (t *Template) Delims(left, right string) *Template {
	t.tpl.Delims(left, right)
	return t
}

// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type or if the name cannot be used syntactically as a function in a template.
// It is legal to overwrite elements of the map. The return value is the template,
// so calls can be chained.
func (t *Template) Funcs(funcMap map[string]interface{}) *Template {
	t.tpl.Funcs(funcMap)
	return t
}

func (t *Template) String() string {
	return t.tpl.DefinedTemplates()
}

//Extension sets template file extension, default is : .html
//only files have the extension will be parsed.
func (t *Template) Extension(ext string) *Template {
	t.ext = ext
	return t
}

// Execute applies the view file associated with t that has the given name
// to the specified data object and return the output.
// If an error occurs executing the template,execution stops.
// A template may be executed safely in parallel.
func (t *Template) Execute(name string, data interface{}) (output []byte, err error) {
	name = strings.Replace(name, "\\", "/", -1)
	if strings.HasSuffix(name, t.ext) {
		name = strings.TrimSuffix(name, t.ext)
	}
	buf := &bytes.Buffer{}
	err = t.tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return
	}
	output = buf.Bytes()
	return
}

// Parse load all view files data and parse it to internal template object.
// Mutiple call of Parse() only the first call worked.
func (t *Template) Parse() (err error) {
	if t.parsed {
		return
	}
	if len(bindata) > 0 {
		err = t.parseFromBinData()
	} else {
		err = t.parseFromDisk()
	}
	if err != nil {
		return
	}
	t.parsed = true
	bindata = nil
	return
}
func (t *Template) parseFromBinData() (err error) {
	for k, v := range bindata {
		html := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", k, string(v))
		_, err = t.tpl.Parse(html)
		if err != nil {
			return
		}
	}
	return
}
func (t *Template) parseFromDisk() (err error) {
	names := []string{}
	err = t.tree(t.rootDir, &names)
	if err != nil {
		return
	}
	var b []byte
	for _, v := range names {
		b, err = ioutil.ReadFile(filepath.Join(t.rootDir, v+t.ext))
		if err != nil {
			return
		}
		html := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", v, string(b))
		_, err = t.tpl.Parse(html)
		if err != nil {
			return
		}
	}
	return
}
func (t *Template) tree(folder string, names *[]string) (err error) {
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
			err = t.tree(v, names)
			if err != nil {
				return err
			}
		} else {
			if !strings.HasSuffix(v, t.ext) {
				continue
			}
			v, err = filepath.Abs(v)
			if err != nil {
				return err
			}
			v0 := strings.Replace(v, "\\", "/", -1)
			v0 = strings.Replace(v0, t.rootDir+"/", "", -1)
			v0 = strings.TrimSuffix(v0, t.ext)
			*names = append(*names, v0)
		}
	}
	if file != nil {
		file.Close()
	}
	return
}
