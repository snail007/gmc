// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gview

import (
	gcore "github.com/snail007/gmc/core"
	ghttputil "github.com/snail007/gmc/internal/util/http"
	"io"
	"strings"
	"sync"
)

type View struct {
	tpl       gcore.Template
	data      map[string]interface{}
	writer    io.Writer
	layout    string
	once      *sync.Once
	onceFn    func()
	lasterr   error
	layoutDir string
}

func New(w io.Writer, tpl gcore.Template) gcore.View {
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
func (this *View) Set(key string, val interface{}) gcore.View {
	this.data[key] = val
	return this
}

// SetMap sets mapped data apply to the template
func (this *View) SetMap(d map[string]interface{}) gcore.View {
	for k, v := range d {
		this.data[k] = v
	}
	return this
}

//Render renders template `tpl` with `data`, and output.
func (this *View) Render(tpl string, data ...map[string]interface{}) gcore.View {
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
	if this.onceFn != nil {
		this.once.Do(this.onceFn)
	}
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
		msg := gcore.ProviderError()().StackError(this.lasterr)
		ghttputil.Stop(this.writer, msg)
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
			msg := gcore.ProviderError()().StackError(this.lasterr)
			ghttputil.Stop(this.writer, msg)
			return
		}
	}
	return
}

// Layout sets the views layout when render template.
func (this *View) Layout(l string) gcore.View {
	this.layout = l
	return this
}

// Stop exit controller method
func (this *View) Stop() {
	ghttputil.Stop(this.writer)
}

// OnRenderOnce injects GPSC data
func (this *View) OnRenderOnce(f func()) gcore.View {
	this.onceFn = f
	return this
}

// SetLayoutDir sets default dir of layout
func (this *View) SetLayoutDir(layoutDir string) {
	layoutDir = strings.Trim(layoutDir, "/")
	layoutDir = strings.Replace(layoutDir, "\\", "/", -1)
	this.layoutDir = layoutDir
}
