package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	gotemplate "text/template"
)

type Template struct {
	rootDir string
	tpl     *gotemplate.Template
	parsed  bool
}

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
	}
	return
}
func (t *Template) String() string {
	return t.tpl.DefinedTemplates()
}
func (t *Template) Execute(name string, data interface{}) (html []byte, err error) {
	name = strings.Replace(name, "\\", "/", -1)
	buf := &bytes.Buffer{}
	err = t.tpl.ExecuteTemplate(buf, "layout/list.html", data)
	html = buf.Bytes()
	return
}
func (t *Template) Parse() (err error) {
	if t.parsed {
		return
	}
	names := []string{}
	err = t.tree(t.rootDir, &names)
	if err != nil {
		return
	}
	for _, v := range names {
		html := fmt.Sprintf("{{define \"%s\"}}\n%s\n{{end}}", v, string(t.getContents(v)))
		_, err = t.tpl.Parse(html)
		if err != nil {
			return
		}
	}
	t.parsed = true
	return
}

func (t *Template) getContents(path string) []byte {
	b, _ := ioutil.ReadFile(filepath.Join(t.rootDir, path))
	return b
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
			v, err = filepath.Abs(v)
			if err != nil {
				return err
			}
			v0 := strings.Replace(v, "\\", "/", -1)
			v0 = strings.Replace(v0, t.rootDir+"/", "", -1)
			*names = append(*names, v0)
		}
	}
	return
}
