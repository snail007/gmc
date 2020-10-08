package gmcview

import (
	gmctemplate "github.com/snail007/gmc/http/template"
	gmchttputil "github.com/snail007/gmc/util/httputil"
	"io"
	"sync"
)

type View struct {
	tpl    *gmctemplate.Template
	data   map[string]interface{}
	writer io.Writer
	layout string
	once   *sync.Once
	onceFn func()
}

func New(w io.Writer, tpl *gmctemplate.Template) *View {
	return &View{
		writer: w,
		tpl:    tpl,
		data:   map[string]interface{}{},
		once:   &sync.Once{},
	}
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
	_, err := this.writer.Write(this.RenderR(tpl, data...))
	if err != nil {
		panic(err)
	}
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
	d, err := this.tpl.Execute(tpl, data0)
	if err != nil {
		panic(err)
	}
	if this.layout != "" {
		data0["GMC_LAYOUT_CONTENT"] = string(d)
		d, err = this.tpl.Execute(this.layout, data0)
		if err != nil {
			panic(err)
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
	this.onceFn=f
	return this
}
