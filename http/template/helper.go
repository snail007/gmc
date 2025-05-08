// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	"embed"
	"encoding/base64"
	"errors"
	gcore "github.com/snail007/gmc/core"
	gctx "github.com/snail007/gmc/module/ctx"
	glog "github.com/snail007/gmc/module/log"
	gcond "github.com/snail007/gmc/util/cond"
	"io/fs"
	"path/filepath"
	"strings"
)

func Init(ctx gcore.Ctx) (tpl gcore.Template, err error) {
	cfg := ctx.Config()
	tpl, err = NewTemplate(ctx, cfg.GetString("template.dir"))
	if err != nil {
		return nil, err
	}
	tpl.Delims(cfg.GetString("template.delimiterleft"),
		cfg.GetString("template.delimiterright"))
	tpl.Extension(cfg.GetString("template.ext"))
	return tpl, nil
}

var ctx = gctx.NewCtx()

func init() {
	ctx.SetLogger(glog.DiscardLogger)
}

func RenderBytes(tpl []byte, data map[string]interface{}) (result []byte, err error) {
	return RenderBytesWithFunc(tpl, data, nil)
}

func RenderBytesWithFunc(tpl []byte, data map[string]interface{}, funcMap map[string]interface{}) (result []byte, err error) {
	t := New()
	t.SetCtx(ctx)
	t.DdisableLogging()
	t.DisableLoadDefaultBinData()
	t.SetBinBytes(map[string][]byte{"tpl": tpl})
	if len(funcMap) > 0 {
		t.Funcs(funcMap)
	}
	err = t.Parse()
	if err != nil {
		return
	}
	return t.Execute("tpl", data)
}

func RenderString(tpl string, data map[string]interface{}) (result string, err error) {
	v, e := RenderBytes([]byte(tpl), data)
	return string(v), e
}

func RenderStringWithFunc(tpl string, data map[string]interface{}, funcMap map[string]interface{}) (result string, err error) {
	r, e := RenderBytesWithFunc([]byte(tpl), data, funcMap)
	return string(r), e
}

type EmbedTemplateFS struct {
	root string
	ext  string
	fs   embed.FS
	tpl  *Template
}

// NewEmbedTemplateFS parse template files from fs embed.FS to tpl *Template,
// it should be called before tpl.Parse(), if the tpl parsed already , nil returned.
func NewEmbedTemplateFS(tpl *Template, fs embed.FS, rootDir string) *EmbedTemplateFS {
	rootDir = strings.Trim(rootDir, "/")
	return &EmbedTemplateFS{
		fs:   fs,
		tpl:  tpl,
		root: rootDir,
		ext:  ".html",
	}
}

func (s *EmbedTemplateFS) SetExt(ext string) *EmbedTemplateFS {
	s.ext = ext
	return s
}

func (s *EmbedTemplateFS) Parse() (err error) {
	if s.tpl.parsed {
		return errors.New("tpl parsed already")
	}
	bindData := map[string]string{}
	err = fs.WalkDir(s.fs, s.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == s.ext {
			b, _ := s.fs.ReadFile(path)
			key := strings.TrimSuffix(path, s.ext)
			trimPrefix := gcond.Cond(s.root == "", "", s.root+"/").String()
			if len(trimPrefix) > 0 {
				key = strings.TrimPrefix(key, trimPrefix)
			}
			bindData[key] = base64.StdEncoding.EncodeToString(b)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.tpl.SetBinBase64(bindData)
	return nil
}
