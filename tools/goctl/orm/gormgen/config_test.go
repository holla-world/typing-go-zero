package gormgen

import (
	"testing"
)

func Test_matchDesc(t *testing.T) {
	// @status(状态): 1-wait(待过期) 2-part(部分使用) 3-all(全部已使用) 4-expired(已过期)
	t.Log(matchDesc("@status(状态)"))
	t.Log(matchKey("@status(状态)"))
}
