// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/zhimaAi/go_tools/tool"
)

func RunJavaScript(mainFunc, batchNo string, params map[string]any) (string, error) {
	vm := goja.New()
	if params == nil {
		params = make(map[string]any)
	}
	err := vm.Set(fmt.Sprintf(`input%s`, batchNo), tool.JsonEncodeNoError(params))
	if err != nil {
		return ``, err
	}
	code := fmt.Sprintf(`%s
	let output%s=JSON.stringify(main(JSON.parse(input%s)))`, mainFunc, batchNo, batchNo)
	_, err = vm.RunString(code)
	if err != nil {
		return ``, err
	}
	return vm.Get(fmt.Sprintf(`output%s`, batchNo)).String(), nil
}
