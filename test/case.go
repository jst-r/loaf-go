package test

import (
	"strings"
	"testing"

	"github.com/jst-r/loaf-go/compiler"
	"github.com/jst-r/loaf-go/vm"
)

const (
	InterpretOk           = vm.InterpretOk
	InterpretErr          = vm.InterpretCompileError
	InterpretRuntimeError = vm.InterpretRuntimeError
)

type Case struct {
	name          string
	code          string
	stdOut        string
	compileErrors []string
}

func NewCase(name string, code string) Case {
	return Case{name: name, code: code}
}

func (c Case) ExpectLines(lines ...string) Case {
	c.stdOut = strings.Join(lines, "\n") + "\n"
	return c
}

func (c Case) ExpectCompileErrors(errs ...string) Case {
	c.compileErrors = errs
	return c
}

func RunCases(cases []Case, t *testing.T) {
	t.Helper()
	for _, c := range cases {
		RunCase(c, t)
	}
}

func RunCase(c Case, t *testing.T) {
	t.Helper()
	t.Run(c.name, func(t *testing.T) {
		prog, errs := compiler.Compile(c.code)
		if c.compileErrors != nil {
			if len(errs) != len(c.compileErrors) {
				t.Logf("Expected %d compile errors, got %d", len(c.compileErrors), len(errs))
				t.Fail()
			}
			for i, err := range errs {
				if err != c.compileErrors[i] {
					t.Logf("Expected error %d to be %s, got %s", i, c.compileErrors[i], err)
					t.Fail()
				}
			}
		} else if len(errs) > 0 {
			t.Log("Compile errors:")
			for _, err := range errs {
				t.Log(err)
			}
			t.Fail()
		}

		vm := vm.New()
		vmOut := strings.Builder{}
		vm.Stdout = &vmOut
		vmErr := strings.Builder{}
		vm.Stderr = &vmErr

		res := vm.Interpret(prog)
		if res != InterpretOk {
			t.Log("Runtime errors:")
			t.Log(vmErr.String())
			t.Fail()
		}

		if c.stdOut != "" {
			if vmOut.String() != c.stdOut {
				t.Logf("\nExpected output:\n%sActual output:\n%s", c.stdOut, vmOut.String())
				t.Fail()
			}
		}

	})
}
