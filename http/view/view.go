package gmcview

import (
	gmcerr "github.com/snail007/gmc/error"
	gmctemplate "github.com/snail007/gmc/http/template"
	gmchttputil "github.com/snail007/gmc/util/http"
	"io"
	"strings"
	"sync"
)

type View struct {
	tpl       *gmctemplate.Template
	data      map[string]interface{}
	writer    io.Writer
	layout    string
	once      *sync.Once
	onceFn    func()
	lasterr   error
	layoutDir string
}

func New(w io.Writer, tpl *gmctemplate.Template) *View {
	return &View{
		writer: w,
		tpl:    tpl,
		data:   map[string]interface{}{},
		once:   &sync.Once{},
	}
}

func (this *View) Err() error {
	return this.lasterr
}

// Set sets data apply to the template
func (this *View) Set(key string, val interface{}) *View {
	this.data[key] = val
	return this
}

// SetMap sets mapped data apply to the template
func (this *View) SetMap(d map[string]interface{}) *View {
	for k, v := range d {
		this.data[k] = v
	}
	return this
}

//Render renders template `tpl` with `data`, and output.
func (this *View) Render(tpl string, data ...map[string]interface{}) *View {
	d := this.RenderR(tpl, data...)
	if this.lasterr != nil {
		return this
	}
	_, this.lasterr = this.writer.Write(d)
	return this
}

//Render renders template `tpl` with `data`, and returns render result.
func (this *View) RenderR(tpl string, data ...map[string]interface{}) (d []byte) {
	// init GPSC
	this.once.Do(this.onceFn)
	data0 := this.data
	for k, v := range this.data {
		data0[k] = v
	}
	if len(data) > 0 {
		for k, v := range data[0] {
			data0[k] = v
		}
	}
	d, this.lasterr = this.tpl.Execute(tpl, data0)
	if this.lasterr != nil {
		gmchttputil.Stop(this.writer, gmcerr.Wrap(this.lasterr).ErrorStack())
		return
	}
	if this.layout != "" {
		data0["GMC_LAYOUT_CONTENT"] = string(d)
		layout := this.layout
		if this.layoutDir != "" {
			layout = this.layoutDir + "/" + this.layout
		}
		d, this.lasterr = this.tpl.Execute(layout, data0)
		if this.lasterr != nil {
			gmchttputil.Stop(this.writer, gmcerr.Wrap(this.lasterr).ErrorStack())
			return
		}
	}
	return
}

// Layout sets the views layout when render template.
func (this *View) Layout(l string) *View {
	this.layout = l
	return this
}

// Stop exit controller method
func (this *View) Stop() {
	gmchttputil.Stop(this.writer)
}

// OnRenderOnce injects GPSC data
func (this *View) OnRenderOnce(f func()) *View {
	this.onceFn = f
	return this
}

// SetLayoutDir sets default dir of layout
func (this *View) SetLayoutDir(layoutDir string) {
	layoutDir = strings.Trim(layoutDir, "/")
	layoutDir = strings.Replace(layoutDir, "\\", "/", -1)
	this.layoutDir = layoutDir
}
