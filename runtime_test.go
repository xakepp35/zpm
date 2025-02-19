package zpm_test

import (
	"fmt"
	"testing"

	"github.com/xakepp35/zpm"
)

type structure struct {
}

func (x *structure) test2() {
	fmt.Println(zpm.RuntimeFuncName(0))
}

func test() {
	fmt.Println(zpm.RuntimeFuncName(0))
	var x structure
	x.test2()
}

func TestRuntimeFunc(t *testing.T) {
	t.SkipNow()
	fmt.Println(zpm.RuntimeFuncName(0))
	test()
}
