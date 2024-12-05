package spec

import (
	"strings"

	"github.com/samber/lo"
)

func IsEnumType(tp DefineStruct) bool {
	if len(tp.Members) == 0 {
		return false
	}
	isEnum := true
	memberType := tp.Members[0].Type.Name()
	for _, member := range tp.Members {
		// 判断是否都是同一个开头
		if !strings.HasPrefix(member.Name, tp.RawName) {
			isEnum = false
			break
		}
		// 判断成员类型是否一致
		if memberType != member.Type.Name() {
			isEnum = false
			break
		}
	}
	return isEnum
}

func ToEnumType(tp DefineStruct) EnumType {
	items := make([]EnumMember, 0)
	var enumGoType string
	for _, member := range tp.Members {
		// 枚举的go的原始类型
		enumGoType = member.Type.Name()
		item := EnumMember{
			Name:    member.Name,
			GoType:  tp.RawName,
			Value:   "",
			Comment: member.Comment,
		}

		switch enumGoType {
		case "int32":
			//  如果是int32，则从tag中获取val的值
			for _, tag := range member.Tags() {
				if tag.Key == "val" {
					item.Value = tag.Name
					break
				}
			}
		case "string":
			// 去掉前缀，然后转下划线
			val := strings.TrimPrefix(member.Name, tp.RawName)
			val = lo.SnakeCase(val)
			item.Value = val
		}

		items = append(items, item)
	}
	return EnumType{
		Name:    tp.RawName,
		GoType:  enumGoType,
		Comment: "",
		Members: items,
	}
}
