package gtemplate

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/util/config"
)

func Init(cfg *gconfig.Config) (tpl gcore.Template, err error) {
	tpl, err = NewTemplate(cfg.GetString("template.dir"))
	if err != nil {
		return nil, err
	}
	tpl.Delims(cfg.GetString("template.delimiterleft"),
		cfg.GetString("template.delimiterright"))
	tpl.Extension(cfg.GetString("template.ext"))
	return tpl, nil
}
