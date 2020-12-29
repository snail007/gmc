package gtemplate

import (
	gcore "github.com/snail007/gmc/core"
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
