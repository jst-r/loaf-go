package test

import "testing"

func TestAssignment(t *testing.T) {
	var cases = []Case{
		NewCase("Plus assign global", `
		var x = 1;
		x += 1;
		print x;
		`).ExpectLines("2"),
		NewCase("Plus assign local", `
		{
			var x = 1;
			x += 1;
			print x;
		}
		`).ExpectLines("2"),
		NewCase("Minus assign global", `
		var x = 1;
		x -= 1;
		print x;
		`).ExpectLines("0"),
		NewCase("Minus assign local", `
		{
			var x = 1;
			x -= 1;
			print x;
		}`).ExpectLines("0"),
	}

	RunCases(cases, t)
}
