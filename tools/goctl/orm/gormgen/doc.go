package gormgen

var doc = `
通过解析表定义里面的Comment信息,实现枚举生成或者自定义映射类型,具体用法如下:
	a.枚举定义:
		1. 整型(tinyint)枚举: 如"@status(状态): 1-wait(待过期) 2-part(部分使用) 3-all(全部已使用) 4-expired(已过期)",其中{status}表示生成go的类型别名,{1}表示枚举值,{part}表示枚举名,{待过期}表示枚举注释
		2. 字符串(char|varchar)枚举: 如"@(通话类型):video_1v1(1v1视频通话) video_match(匹配视频通话)",表述和整型枚举基本一致,去掉了数字部分
	b.引用枚举: 支持引用其他表已定义好的枚举,如"@[CRoomStatus](房间开播状态)",[]:表示引用,CRoomStatus:表示引用对象的go类型名称,(房间开播状态):表示注释
	c.指定类型: 在Comment只要包含以下(float32|float64|decimal)关键词,可将字段映射为指定类型,目前支持:float32(原生类型);float64(原生类型);decimal(decimal包,适合高精度计算)
`