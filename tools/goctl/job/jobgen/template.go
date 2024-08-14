package jobgen

import (
	_ "embed"
)

// 预设模板文件名
const (
	category                   = "job"
	handlerTemplateFile        = "handler.tpl"
	routesTemplateFile         = "routes.tpl"
	routesAdditionTemplateFile = "route-addition.tpl"
)

//go:embed handler.tpl
var handlerTemplate string

//go:embed routes.tpl
var routesTemplate string

//go:embed route-addition.tpl
var routesAdditionTemplate string
