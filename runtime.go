package zpm

import (
	"runtime"
	"strings"
)

func RuntimeFuncName(skip int) string {
	fn := runtimeFunc(skip + 1)
	res := fn.Name()
	lastIndex := strings.LastIndex(res, "/")
	res = res[lastIndex+1:]
	return res
}

func runtimeFunc(skip int) *runtime.Func {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	}
	return runtime.FuncForPC(pc)
}
