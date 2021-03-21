package test

import (
	"runtime"
	"strings"
	"time"
)

var testRunID = time.Now().Unix()

func AbsoluteProjectDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return strings.Split(filename, "/todosync/test/")[0]
}
