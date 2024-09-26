package gormgen

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/samber/lo"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	// mysql type to golang type
	dataTypeMap = map[string]func(columnType gorm.ColumnType) (dataType string){
		"bigint": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"int": func(columnType gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"decimal": func(columnType gorm.ColumnType) (dataType string) {
			return "float32"
		},
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			return "int32"
		},
	}
)

type Enum struct {
	GenType string
	Name    string
	Value   string
	Comment string
}

type EnumField struct {
	Name       string
	NativeType string
	GenType    string
	Comment    string
	Enums      []Enum
}

type Spec struct {
	TableName      string      `json:"table_name,omitempty"`
	Shards         int         `json:"shards,omitempty"`
	ModelName      string      `json:"model,omitempty"`
	ShardKeyName   string      `json:"shard_key,omitempty"`
	ShardKeyGoType string      `json:"shard_key_go_type,omitempty"`
	EnumFields     []EnumField `json:"enum_fields"`
}

type GenerateSpec struct {
	TableSpec []*Spec
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

// 全命名格式：
// @status(状态): 1-wait(待过期) 2-part(部分使用) 3-all(全部已使用) 4-expired(已过期)
// 缺省名称和部分枚举描述
// @(状态): 1-wait 2-part 3-all(全部已使用) 4-expired(已过期)
func ParseEnum(modelName string, gf gen.Field) (ef EnumField, ok bool) {
	if gf == nil {
		return
	}
	comment := gf.ColumnComment
	if !strings.HasPrefix(comment, "@") {
		return
	}
	sp := strings.SplitN(comment, ":", 2)
	if strings.Contains(comment, ":") && len(sp) > 0 {
		// 解析头部描述
		// @status(状态)
		head := sp[0]
		desc := matchDesc(head)
		// 必须要有描述
		if desc == "" {
			return
		}
		ok = true
		ef.Comment = desc
		ef.Name = matchKey(head)
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

func matchDesc(text string) string {
	re := regexp.MustCompile(`[{}()\[\]](.*?)[{}()\[\]]`)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		return match[1] // 输出第一个括号中的内容
	}
	return ""
}

func matchKey(text string) string {
	// 正则表达式，匹配仅包含字母和数字的单词
	pattern := `\b[a-zA-Z0-9]+\b`
	re := regexp.MustCompile(pattern)
	return re.FindString(text)
}
