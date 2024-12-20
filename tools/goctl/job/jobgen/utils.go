package jobgen

import (
	"bytes"
	cron2 "github.com/lnquy/cron"
	"github.com/robfig/cron/v3"
	"regexp"
	"strings"
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
	description, err := parser.ToDescription(expr, cron2.Locale_en)
	if err != nil {
		return "", err
	}
	return description, nil
}

func toCamelCase(input string) string {
	// 正则匹配分隔符（- 或 _）
	re := regexp.MustCompile(`[-_]`)
	// 按分隔符分割字符串
	parts := re.Split(input, -1)
	// 首字母大写
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	// 拼接成单个字符串
	return strings.Join(parts, "")
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	// 获取首字符并转换为小写，拼接剩余部分
	return strings.ToLower(string(s[0])) + s[1:]
}