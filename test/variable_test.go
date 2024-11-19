package test

import (
	"testing"
)

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
}

func TestVariables(t *testing.T) {
	for _, c := range cases {
		RunCase(c, t)
	}
}
