package gormgen

import (
	"text/template"
)

// 分表方法模板
func shardingTpl() string {
	return `// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
package query

import (
	"strconv"
)

type Numeric interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

// sharding function that uses generics to accept different numeric types
func ShardingNo[T Numeric](tableNum int, key T) int {
	return int(key) % tableNum
}

// sharding function that uses generics to accept different numeric types
func Sharding[T Numeric](tableName string, tableNum int, key T) string {
	mod := ShardingNo(tableNum, key)
	return tableName + "_" + strconv.Itoa(mod)
}`
}

// 分表方法模板
func tbMethodTpl() (*template.Template, error) {
	tpl := `
func (q *Query) {{.ModelName}}TB({{.ShardKeyName}} {{.ShardKeyGoType}}) *{{.Instantiation}} {
	t := Sharding("{{.Name}}", {{.Shards}}, {{.ShardKeyName}})
	return q.{{.Instantiation}}.Table(t)
}

func (q *Query) {{.ModelName}}AssignTBName(name string) *{{.Instantiation}} {
	return q.{{.Instantiation}}.Table(name)
}`
	return template.New("meth").Parse(tpl)
}

// 分表元数据
func shardingMetaTpl() (*template.Template, error) {
	tpl := `
// 分表数量
func (m {{.ModelName}}) Shards() int64 {
	return {{.Shards}}
}

// 分表序号，在第几个表
func (m {{.ModelName}}) ShardSeq({{.ShardKeyNameCamel}} int64) int64 {
	return {{.ShardKeyNameCamel}} % m.Shards()
}

// 分表key
func (m {{.ModelName}}) ShardKey() string {
	return "{{.ShardKeyName}}"
}

// 分表分组,给定一组id,按所在分表分组后的结果
func (m {{.ModelName}}) ShardGroup({{.ShardKeyNameCamel}}s []int64) map[int64][]int64 {
	g := map[int64][]int64{}
	for _, id := range {{.ShardKeyNameCamel}}s {
		key := id % m.Shards()
		g[key] = append(g[key], id)
	}
	return g
}
`
	return template.New("meta").Parse(tpl)
}
