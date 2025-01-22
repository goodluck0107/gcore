package utils_test

import (
	"gitee.com/monobytes/gcore/glog/zap/internal/utils"
	"testing"
)

func TestAddslashes(t *testing.T) {
	str1 := "abc\\mas"
	t.Log(str1)
	t.Log(utils.Addslashes(str1))

	str2 := "abc\"mas"
	t.Log(str2)
	t.Log(utils.Addslashes(str2))

	str3 := "abc'mas"
	t.Log(str3)
	t.Log(utils.Addslashes(str3))

	str4 := "abc\nmas"
	t.Log(str4)
	t.Log(utils.Addslashes(str4))

	str5 := "abc\tmas"
	t.Log(str5)
	t.Log(utils.Addslashes(str5))
}
