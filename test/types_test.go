package test

import (
	"testing"
)

func TestTypes_instanceof(t *testing.T) {
	tc := GivenContextWith(t, "../lib/types.lua", "../lib/Semver.lua", "../lib/Version.lua")

	cases := []struct {
		givenObj   string
		givenClass string
		expected   bool
	}{
		{`nil`, `nil`, true},
		{`"abc"`, `string`, true},
		{`Semver:new("1.2.3")`, `Semver`, true},
		{`Semver:new("1")`, `Semver`, false},
		{`Semver:new("1.2.3")`, `nil`, false},
		{`Semver:new("1")`, `nil`, true},
		{`Version:new("")`, `Version`, false},
		{`Version:new("1.2.3")`, `Version`, true},
		{`Version:new("")`, `nil`, true},
		{`Version:new("1.2.3")`, `nil`, false},
	}

	for _, c := range cases {
		t.Run(c.givenObj+"_"+c.givenClass, func(t *testing.T) {
			tc.ShouldEvaluateTo(t, `return t.instanceof(`+c.givenObj+`, `+c.givenClass+`)`, c.expected)
		})
	}
}
