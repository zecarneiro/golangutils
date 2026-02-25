package logic

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"golangutils/pkg/logger"
	"golangutils/pkg/testsuite"
)

func Ternary[T any](condition bool, resultTrue, resultFalse T) T {
	if condition {
		return resultTrue
	}
	return resultFalse
}

func Exit(code int) {
	if testsuite.IsRunningTests {
		logger.Info(fmt.Sprintf("Detect Exit with code: %d", code))
	} else {
		os.Exit(code)
	}
}

func ProcessError(err error) {
	if err != nil {
		logger.Error(err)
		Exit(1)
	}
}

func CaptureOutput(f any, args ...any) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	v := reflect.ValueOf(f)
	relArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		relArgs[i] = reflect.ValueOf(arg)
	}
	v.Call(relArgs)
	w.Close()
	os.Stdout = old
	out := <-outC
	return out, nil
}

// Example: map[string]interface{}{"FUNC_NAME": FUNC, "FUNC_NAME_1": FUNC_1, ....}
type TaskFunc func()

func FuncCall[T interface{}](caller interface{}, params ...interface{}) (T, error) {
	var in []reflect.Value = []reflect.Value{}
	var result T
	var err error
	funcRef := reflect.ValueOf(caller)
	if len(params) > 0 {
		if len(params) != funcRef.Type().NumIn() {
			err = errors.New("The number of params is out of index.")
		}
	}
	if err == nil {
		in = make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		res := funcRef.Call(in)
		if res != nil {
			result = res[0].Interface().(T)
		}
	} else {
		result = *new(T)
	}
	return result, err
}
