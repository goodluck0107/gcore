package gvalidate_test

import (
	"github.com/goodluck0107/gcore/gutils/gvalidate"
	"testing"
)

// "XXXX-XXXXXXX"
// "XXXX-XXXXXXXX"
// "XXX-XXXXXXX"
// "XXX-XXXXXXXX"
// "XXXXXXX"
// "XXXXXXXX"
func TestIsTelephone(t *testing.T) {
	t.Log(gvalidate.IsTelephone("0285-5554540"))
	t.Log(gvalidate.IsTelephone("0285-55545401"))
	t.Log(gvalidate.IsTelephone("028-5554540"))
	t.Log(gvalidate.IsTelephone("028-55545401"))
	t.Log(gvalidate.IsTelephone("5554540"))
	t.Log(gvalidate.IsTelephone("55545401"))
}

func TestIsEmail(t *testing.T) {
	t.Log(gvalidate.IsEmail("yuebanfuxiao@gmail.com"))
	t.Log(gvalidate.IsEmail("yuebanfuxiao"))
}

func TestIsAccount(t *testing.T) {
	t.Log(gvalidate.IsAccount("0abc", 4, 8))
	t.Log(gvalidate.IsAccount("abc0", 4, 8))
	t.Log(gvalidate.IsAccount("ab0", 4, 8))
	t.Log(gvalidate.IsAccount("ab.cd", 4, 8))
	t.Log(gvalidate.IsAccount("ab_cd", 4, 8))
	t.Log(gvalidate.IsAccount("ab cd", 4, 8))
	t.Log(gvalidate.IsAccount("ab-cd", 4, 8))
	t.Log(gvalidate.IsAccount("abcdefghi", 4, 8))
	t.Log(gvalidate.IsAccount("abcdefghif", 4, 8))
}

func TestIsUrl(t *testing.T) {
	t.Log(gvalidate.IsUrl("http://www.baidu.com"))
	t.Log(gvalidate.IsUrl("HTTP://WWW.BAIDU.COM"))
	t.Log(gvalidate.IsUrl("HTTP://a.b"))
	t.Log(gvalidate.IsUrl("HTTPs://a.b"))
}

func TestIsDigit(t *testing.T) {
	t.Log(gvalidate.IsDigit("11"))
	t.Log(gvalidate.IsDigit("11."))
	t.Log(gvalidate.IsDigit("11.1"))
	t.Log(gvalidate.IsDigit("011.1"))
	t.Log(gvalidate.IsDigit("aa.1"))
}

func TestIn(t *testing.T) {
	t.Log(gvalidate.In("a", []string{"a", "b", "c"}))
	t.Log(gvalidate.In("d", []string{"a", "b", "c"}))
	t.Log(gvalidate.In([]string{"a", "b", "c"}, []string{"a", "b", "c"}))
	t.Log(gvalidate.In([]string{"a", "b", "c"}, []string{"d", "f", "g"}))
}

func TestIsIdCard(t *testing.T) {
	t.Log(gvalidate.IsIdCard("512301195011260279"))
}
