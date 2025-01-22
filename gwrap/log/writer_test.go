package log_test

import (
	"gitee.com/monobytes/gcore/gutils/grand"
	"gitee.com/monobytes/gcore/gwrap/log"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	str := grand.Letters(log.KB) + "\n"

	w := log.NewWriter(
		log.WithFileMaxSize(2*log.KB),
		log.WithFileRotate(log.FileRotateByMinute),
		log.WithCompress(false),
	)

	for i := 0; i < 10; i++ {
		if _, err := w.Write([]byte(str)); err != nil {
			t.Fatal(err)
		}
	}
}
