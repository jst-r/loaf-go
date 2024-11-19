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
	name   string
	code   string
	stdOut string
}

func NewCase(name string, code string) Case {
	return Case{name: name, code: code}
}

func (c Case) ExpectLines(lines ...string) Case {
	c.stdOut = strings.Join(lines, "\n") + "\n"
	return c
}

var cases = []Case{
	NewCase("Globals", `
	var x = 1;
	print x;
	x = x + 1;
	print x;
	`).ExpectLines("1", "2"),
	NewCase("Locals", `
	{
		var x = 1;
		var y = 2;
		print x;
		print y;
		y = 3;
		print y;
		print x + y;
	}
	`).ExpectLines("1", "2", "3", "4"),
	NewCase("Local shadowing", `
	var x = 1;
	{
		var x = 2;
		print x;
	}
	print x;
	`).ExpectLines("2", "1"),
}

func TestIntegration(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			prog, errs := compiler.Compile(c.code)
			if len(errs) > 0 {
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
}
