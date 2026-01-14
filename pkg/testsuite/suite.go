package testsuite

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Suite define a estrutura base para um grupo de testes com estado compartilhado
type Suite struct {
	BeforeEach func(*testing.T)
	AfterEach  func(*testing.T)
}

func New() *Suite {
	IsRunningTests = true
	return &Suite{}
}

func NewWithAll(beforeEach func(*testing.T), afterEach func(*testing.T)) *Suite {
	IsRunningTests = true
	return &Suite{
		BeforeEach: beforeEach,
		AfterEach:  afterEach,
	}
}

func NewWithBeforeEach(beforeEach func(*testing.T)) *Suite {
	IsRunningTests = true
	return &Suite{
		BeforeEach: beforeEach,
	}
}

func NewWithAfterEach(afterEach func(*testing.T)) *Suite {
	IsRunningTests = true
	return &Suite{
		AfterEach: afterEach,
	}
}

func getFunctionName(i interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	if fn == nil {
		return "Test"
	}
	parts := strings.Split(fn.Name(), ".")
	return parts[len(parts)-1]
}

func resolveParamName(p any, index int) string {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Struct {
		for _, name := range []string{"Name", "Description"} {
			field := v.FieldByName(name)
			if field.IsValid() && field.Kind() == reflect.String && field.String() != "" {
				return field.String()
			}
		}
	}
	return fmt.Sprintf("case_%d", index+1)
}

func RunN(t *testing.T, s *Suite, name string, fn func(*testing.T)) {
	t.Run(name, func(st *testing.T) {
		if s.BeforeEach != nil {
			s.BeforeEach(st)
		}
		defer func() {
			if s.AfterEach != nil {
				s.AfterEach(st)
			}
		}()
		fn(st)
	})
}

func Run(t *testing.T, s *Suite, fn func(*testing.T)) {
	name := getFunctionName(fn)
	RunN(t, s, name, fn)
}

func RunTestCasesN[TData any](t *testing.T, s *Suite, name string, dataProviders []TData, fn func(*testing.T, TData)) {
	for i, dp := range dataProviders {
		caseName := fmt.Sprintf("%s/%s", name, resolveParamName(dp, i))
		t.Run(caseName, func(st *testing.T) {
			if s.BeforeEach != nil {
				s.BeforeEach(st)
			}
			defer func() {
				if s.AfterEach != nil {
					s.AfterEach(st)
				}
			}()
			fn(st, dp)
		})
	}
}

func RunTestCases[TData any](t *testing.T, s *Suite, dataProviders []TData, fn func(*testing.T, TData)) {
	name := getFunctionName(fn)
	RunTestCasesN(t, s, name, dataProviders, fn)
}
