package gormgen

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func Gen(src, pkg string) error {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/repo/" + pkg,
		ModelPkgPath:      "internal/repo/model",
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})
	db, err := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{
			src,
		},
	}))
	if err != nil {
		return err
	}
	specs, err := parseTables(src)
	if err != nil {
		return err
	}
	cfg := GenerateSpec{
		TableSpec: specs,
	}

	g.UseDB(db)
	g.WithDataTypeMap(dataTypeMap)
	for _, v := range cfg.TableSpec {
		g.ApplyBasic(g.GenerateModelAs(v.TableName, v.ModelName))
	}

	g.Execute()
	err = insertTbMethod(cfg, g)
	if err != nil {
		return err
	}
	err = insertModelMeth(cfg, g)
	if err != nil {
		return err
	}
	err = createShardingFile(g, pkg)
	if err != nil {
		return err
	}
	printDDL(cfg)
	return nil
}

// parseTables
func parseTables(file string) ([]Spec, error) {
	ddl, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer ddl.Close()

	specs := make([]Spec, 0, 10)
	scanner := bufio.NewScanner(ddl)
	for scanner.Scan() {
		line := scanner.Text()
		ts := strings.ReplaceAll(line, " ", "")
		ls := strings.ToLower(ts)
		if !strings.HasPrefix(ls, "#@meta") && !strings.HasPrefix(ls, "createtable") {
			continue
		}
		if strings.HasSuffix(ls, ";") {
			continue
		}
		spec := Spec{}
		if strings.HasPrefix(ls, "createtable") {
			spec.TableName = parseTableName(line)
			spec.ModelName = ToCamelCase(spec.TableName)
			specs = append(specs, spec)
			continue
		}
		ts = strings.TrimPrefix(ts, "#@meta")
		err := json.Unmarshal([]byte(ts), &spec)
		if err != nil {
			return nil, fmt.Errorf("解析meta错误: %s\n  line:[%s]", err, line)
		}
		if spec.Shards > 1 && spec.ShardKeyName == "" {
			return nil, fmt.Errorf("分表必须指定shards和shard_key_name\n   line:[%s]", line)
		}
		// 扫描表名
		for i := 0; i < 10; i++ {
			ok := scanner.Scan()
			if !ok {
				return nil, fmt.Errorf("解析meta错误，meta之后应该紧跟建表语句\n   line:[%s]", line)
			}
			line = scanner.Text()
			tn := parseTableName(line)
			if tn == "" {
				continue
			}
			spec.TableName = tn
			break
		}
		if spec.TableName == "" {
			return nil, fmt.Errorf("解析meta错误，meta之后应该紧跟建表语句\n   line:[%s]", line)
		}
		// 给默认值
		if spec.ModelName == "" {
			spec.ModelName = ToCamelCase(spec.TableName)
		}
		// 给默认值
		if spec.ShardKeyGoType == "" {
			spec.ShardKeyGoType = "int64"
		}
		specs = append(specs, spec)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return specs, nil
}

func parseTableName(origin string) string {
	ls := strings.ToLower(origin)
	if !strings.HasPrefix(ls, "create table") {
		return ""
	}
	// 解析表名
	rs := strings.Replace(ls, "create table", "", 1)
	// 去掉空格
	rs = strings.ReplaceAll(rs, " ", "")
	// 去掉反引号
	rs = strings.ReplaceAll(rs, "`", "")
	// 去掉单引号
	rs = strings.ReplaceAll(rs, "'", "")
	// .号分割
	ss := strings.Split(rs, ".")
	if len(ss) == 2 {
		return ss[1]
	} else {
		return ss[0]
	}
}

func printDDL(cfg GenerateSpec) {
	// CREATE TABLE post_likes_1 LIKE post_likes;
	for _, v := range cfg.TableSpec {
		if v.Shards <= 1 {
			continue
		}
		for i := 0; i < v.Shards; i++ {
			fmt.Printf("CREATE TABLE %s_%d LIKE %s;\n", v.TableName, i, v.TableName)
		}
	}
}

func insertTbMethod(cfg GenerateSpec, g *gen.Generator) error {
	genFile, err := os.Open(g.OutFile)
	if err != nil {
		return err
	}
	defer genFile.Close()

	var fileContent strings.Builder

	scanner := bufio.NewScanner(genFile)
	regexps := buildReplaceRegexp(cfg)

	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range regexps {
			line = v.regex.ReplaceAllString(line, v.unexported)
		}
		fileContent.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	tpl, err := tbMethodTpl()
	if err != nil {
		return err
	}
	for _, v := range cfg.TableSpec {
		if v.Shards <= 1 {
			continue
		}

		var rendered bytes.Buffer
		err = tpl.Execute(&rendered, map[string]interface{}{
			"Name":              v.TableName,
			"ModelName":         v.ModelName,
			"Shards":            v.Shards,
			"ShardKeyGoType":    v.ShardKeyGoType,
			"ShardKeyName":      v.ShardKeyName,
			"ShardKeyNameCamel": convertSnakeToCamel(v.ShardKeyName),
			"Instantiation":     v.UnExportName(),
		})

		if err != nil {
			return err
		}
		s := rendered.String()
		fileContent.WriteString(s + "\n")
	}
	outputFile, err := os.Create(g.OutFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(fileContent.String())
	return err
}

func insertModelMeth(cfg GenerateSpec, g *gen.Generator) error {
	tpl, err := shardingMetaTpl()
	if err != nil {
		return err
	}
	for _, v := range cfg.TableSpec {
		if v.Shards <= 1 {
			continue
		}
		var rendered bytes.Buffer
		err = tpl.Execute(&rendered, map[string]interface{}{
			"Name":              v.TableName,
			"ModelName":         v.ModelName,
			"Shards":            v.Shards,
			"ShardKeyGoType":    v.ShardKeyGoType,
			"ShardKeyName":      v.ShardKeyName,
			"ShardKeyNameCamel": convertSnakeToCamel(v.ShardKeyName),
			"Instantiation":     v.UnExportName(),
		})
		if err != nil {
			return err
		}
		appendText := rendered.String()

		filePath := path.Join(g.ModelPkgPath, fmt.Sprintf("%s.gen.go", v.TableName))
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.WriteString(appendText); err != nil {
			return err
		}
	}

	return nil
}

type replaceRegexp struct {
	regex      *regexp.Regexp
	exported   string
	unexported string
}

// buildReplaceRegexp build regexp for insertTbMethod
func buildReplaceRegexp(c GenerateSpec) []*replaceRegexp {
	rrs := make([]*replaceRegexp, 0, len(c.TableSpec))
	for _, v := range c.TableSpec {
		if v.Shards <= 1 {
			continue
		}
		regex := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, v.ModelName))
		rrs = append(rrs, &replaceRegexp{
			regex:      regex,
			exported:   v.ModelName,
			unexported: v.UnExportName(),
		})
	}
	return rrs
}

func createShardingFile(g *gen.Generator, pkg string) error {
	filename := fmt.Sprintf("%s/sharding.gen.go", g.OutPath)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		content := shardingTpl()
		// 解析模板
		parse, err := template.New("shard").Parse(content)
		if err != nil {
			return err
		}
		err = parse.Execute(file, map[string]interface{}{
			"Package": filepath.Base(pkg),
		})
		if err != nil {
			return err
		}

		fmt.Println("create sharding file success")
		return nil
	} else {
		fmt.Println("sharding file already exists")
		return nil
	}
}

func convertSnakeToCamel(input string) string {
	var result strings.Builder
	convertNext := false

	for i, v := range input {
		if v == '_' {
			convertNext = true
		} else {
			if convertNext {
				result.WriteRune(unicode.ToUpper(v))
				convertNext = false
			} else {
				if i == 0 {
					result.WriteRune(unicode.ToLower(v))
				} else {
					result.WriteRune(v)
				}
			}
		}
	}
	return result.String()
}

// ToCamelCase 将任意字符串转换为 CamelCase 形式
func ToCamelCase(s string) string {
	// 分割字符串为单词列表
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// 转换每个单词为首字母大写
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	// 拼接所有单词为一个字符串
	return strings.Join(words, "")
}
