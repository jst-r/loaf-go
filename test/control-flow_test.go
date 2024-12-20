package test

import "testing"

func TestIfElse(t *testing.T) {
	var cases = []Case{
		NewCase("If", `
	var x = 1;
	if (x == 1) {
		print "x is 1";
	}`).ExpectLines("x is 1").ExpectStackSize(0),

		NewCase("If else", `
	var x = 1;
	if (x == 2) {
		print "x is 1";
	} else {
		print "x is not 1";
	}`).ExpectLines("x is not 1").ExpectStackSize(0),

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
	`).ExpectLines("2", "3").ExpectStackSize(0),
	}
	RunCases(cases, t)
}

func TestLogicalOps(t *testing.T) {
	c := NewCase("logical ops", `
	print true and false; // false
	print false and true; // false
	print true and 1;     // 1
	print nil and true;	  // nil
	print true or false;  // true
	print false or true;  // true
	print nil or true;    // true
	print 1 or nil;       // 1
	print false or 1;     // 1
	`).ExpectLines("false", "false", "1", "nil", "true", "true", "true", "1", "1")

	RunCase(c, t)
}

func TestWhileLoop(t *testing.T) {
	c := NewCase("", `
	var i = 1;
	while i <= 5 {
		print i;
		i = i + 1;
	}
	`).ExpectLines("1", "2", "3", "4", "5")

	RunCase(c, t)
}
