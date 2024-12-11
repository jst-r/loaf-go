package test

import "testing"

func TestControlFlow(t *testing.T) {
	var cases = []Case{
		NewCase("If", `
	var x = 1;
	if (x == 1) {
		print "x is 1";
	}`).ExpectLines("x is 1"),

		NewCase("If else", `
	var x = 1;
	if (x == 2) {
		print "x is 1";
	} else {
		print "x is not 1";
	}`).ExpectLines("x is not 1"),

		NewCase("Locals", `
	var x = 1;
	if (true) {
		x = 2;
	}
	print x;
	if (false) {
	} else {
		x = 3;
	}
	print x;
	`).ExpectLines("2", "3"),
	}
	RunCases(cases, t)
}
