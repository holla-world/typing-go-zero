package gormgen

import (
	"unicode"

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

type Spec struct {
	TableName      string `json:"table_name,omitempty"`
	Shards         int    `json:"shards,omitempty"`
	ModelName      string `json:"model,omitempty"`
	ShardKeyName   string `json:"shard_key,omitempty"`
	ShardKeyGoType string `json:"shard_key_go_type,omitempty"`
}

type GenerateSpec struct {
	TableSpec []Spec
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
