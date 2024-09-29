package gormgen

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/samber/lo"
	"gorm.io/gen"
)

type Enum struct {
	GenType string
	Name    string
	Value   string
	Comment string
}

type EnumField struct {
	Name       string // 枚举名称
	OriginKey  string //
	NativeType string // 原生类型
	GenType    string // 生成类型
	IsCite     bool   // 是否引用
	Comment    string // 注释
	Enums      []Enum // 枚举值
}

// 匹配{}或者()里面的内容
func matchDesc(text string) string {
	re := regexp.MustCompile(`[{}()](.*?)[{}()]`)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// 匹配[]里面的内容
func matchCite(text string) string {
	re := regexp.MustCompile(`[\[\]](.*?)[\[\]]`)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// 匹配单词
func matchKey(text string) string {
	// 正则表达式，匹配仅包含字母和数字的单词
	pattern := `\b[a-zA-Z0-9_-]+\b`
	re := regexp.MustCompile(pattern)
	return re.FindString(text)
}
func lowercaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func (s *Spec) UnExportName() string {
	return lowercaseFirst(s.ModelName)
}

// 支持以下几种格式
// 全命名格式 -> @status(状态): 1-wait(待过期) 2-part(部分使用) 3-all(全部已使用) 4-expired(已过期)
// 缺省名称和部分枚举描述格式 -> @(状态): 1-wait 2-part 3-all(全部已使用) 4-expired(已过期)
// 引用其他model的枚举格式 -> @[CommonStatus](状态)
func ParseEnum(modelName string, gf gen.Field) (ef EnumField, ok bool) {
	if gf == nil {
		return
	}
	comment := gf.ColumnComment
	if !strings.HasPrefix(comment, "@") {
		return
	}

	sp := strings.SplitN(comment, ":", 2)
	if len(sp) > 0 {
		// 解析头部描述
		// @status[CommonStatus](状态)
		head := sp[0]
		desc := matchDesc(head) // desc=状态
		// 必须要有描述
		if desc == "" {
			return
		}
		ok = true
		ef.Comment = desc

		key := matchKey(head) // key=status
		ef.OriginKey = key
		// 是否为引用,引用无需解析枚举值
		cite := matchCite(head) // cite=CommonStatus
		if cite != "" {
			ef.GenType = cite
			ef.IsCite = true
			return
		}

		ef.Name = key
		if ef.Name == "" {
			ef.Name = gf.Name
		}
		// 转驼峰
		ef.Name = lo.PascalCase(modelName) + lo.PascalCase(ef.Name)
		ef.GenType = ef.Name
		ef.NativeType = gf.Type
		// 解析尾部枚举值
		if len(sp) > 1 {
			enums := strings.Split(sp[1], " ")
			result := make([]Enum, 0, len(enums))
			for _, enum := range enums {
				// 1-wait(待过期)
				enum = strings.TrimSpace(enum)
				mixed := strings.Split(enum, "-")
				if len(mixed) < 1 {
					continue
				}

				eval := mixed[0]
				if eval == "" {
					continue
				}
				if gf.Type == "string" {
					eval = fmt.Sprintf(`"%s"`, eval)
				}

				name := matchKey(mixed[1])
				if name == "" {
					continue
				}
				name = ef.Name + lo.PascalCase(name)

				result = append(result, Enum{
					Name:    name,
					Value:   eval,
					Comment: matchDesc(enum),
					GenType: ef.GenType,
				})
			}
			ef.Enums = result
		}
	}
	return
}
