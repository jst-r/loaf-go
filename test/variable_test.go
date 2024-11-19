package test

import (
	"testing"
)

func TestVariables(t *testing.T) {
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
			{
				var x = 3;
				print x;
			}
			print x;
		}
		print x;
		`).ExpectLines("3", "2", "1"),
		NewCase("Shadow with use", `
		{
			var x = 1;
			{
				var x = x + 1;
				print x;
			}
			print x;
		}
		`).ExpectCompileErrors("[line 5]error at x: can not read local variable in its own initializer"),
	}
	RunCases(cases, t)
}
