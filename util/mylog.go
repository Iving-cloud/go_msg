package util

import (
	"fmt"
	"time"
)

type MyLog struct {
	buff string
}

func (this *MyLog) Printf(format string, v ...interface{}) {
	this.buff = fmt.Sprintf("mylog ", time.Now(), "hahahh")
}
