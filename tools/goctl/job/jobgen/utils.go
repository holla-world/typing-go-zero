package jobgen

import (
	"bytes"
	cron2 "github.com/lnquy/cron"
	"github.com/robfig/cron/v3"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type fileGenConfig struct {
	dir             string
	subdir          string
	filename        string
	templateName    string
	category        string
	templateFile    string
	builtinTemplate string
	data            any
}

func genFile(c fileGenConfig) error {
	fp, created, err := util.MaybeCreateFile(c.dir, c.subdir, c.filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var text string
	if len(c.category) == 0 || len(c.templateFile) == 0 {
		text = c.builtinTemplate
	} else {
		text, err = pathx.LoadTemplate(c.category, c.templateFile, c.builtinTemplate)
		if err != nil {
			return err
		}
	}

	t := template.Must(template.New(c.templateName).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, c.data)
	if err != nil {
		return err
	}

	code := golang.FormatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}

// ValidCronExpression checks if a cron expression is valid
func ValidCronExpression(expr string) error {
	_, err := cron.ParseStandard(expr)
	return err
}

func CronExpressionTxt(expr string) (string, error) {
	// 创建一个新的 CronParser
	parser, _ := cron2.NewDescriptor()

	// 将 cron 表达式转换为自然语言描述
	description, err := parser.ToDescription(expr, cron2.Locale_zh_CN)
	if err != nil {
		return "", err
	}
	return description, nil
}
