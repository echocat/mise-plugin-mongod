package test

import (
	"testing"
)

func TestVersion_new(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Version.lua")

	cases := []struct {
		input       string
		expected    any
		expectedErr string
	}{
		{`"1"`, []any{float64(1)}, ""},
		{`"1.2"`, []any{float64(1), float64(2)}, ""},
		{`"1.2.3"`, []any{float64(1), float64(2), float64(3)}, ""},
		{`"1.2.0"`, []any{float64(1), float64(2), float64(0)}, ""},
		{`"0.0.0"`, []any{float64(0), float64(0), float64(0)}, ""},
		{`{}`, nil, ""},
		{`{1}`, []any{float64(1)}, ""},
		{`{1,2}`, []any{float64(1), float64(2)}, ""},
		{`{1,2,3}`, []any{float64(1), float64(2), float64(3)}, ""},
		{`1`, nil, "requires a string or table(array) to create a version from; but got number"},
		{`"a"`, nil, ""},
		{`"1.b"`, nil, ""},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if expectedErr := c.expectedErr; expectedErr == "" {
				tc.ShouldEvaluateTo(t, `return t:new(`+c.input+`)`, c.expected)
			} else {
				tc.ShouldEvaluateToError(t, `return t:new(`+c.input+`)`, expectedErr)
			}
		})
	}
}

func TestVersion_cmp(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Version.lua")

	cases := []struct {
		a, b     string
		expected float64
	}{
		{"1", "1", 0},
		{"2", "1", 1},
		{"1", "2", -1},

		{"1.1", "1.1", 0},
		{"1.0", "1.0", 0},
		{"1.0", "1", 0},
		{"1", "1.0", 0},
		{"1.2", "1.1", 1},
		{"1.1", "1.2", -1},
		{"2.1", "1.1", 1},
		{"1.1", "2.1", -1},
		{"1.1", "1.0", 1},
		{"1", "1.1", -1},
	}

	t.Run("by_string", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.a+"_"+c.b, func(t *testing.T) {
				tc.ShouldEvaluateTo(t, `return t.cmp("`+c.a+`","`+c.b+`")`, c.expected)
			})
		}
	})
	t.Run("instance_by_string", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.a+"_"+c.b, func(t *testing.T) {
				tc.ShouldEvaluateTo(t, `return t:new("`+c.a+`"):cmp("`+c.b+`")`, c.expected)
			})
		}
	})
	t.Run("by_both_instances", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.a+"_"+c.b, func(t *testing.T) {
				tc.ShouldEvaluateTo(t, `return t.cmp(t:new("`+c.a+`"),t:new("`+c.b+`"))`, c.expected)
			})
		}
	})
	t.Run("by_a_instance_b_string", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.a+"_"+c.b, func(t *testing.T) {
				tc.ShouldEvaluateTo(t, `return t.cmp(t:new("`+c.a+`"),"`+c.b+`")`, c.expected)
			})
		}
	})
	t.Run("by_a_string_b_instance", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.a+"_"+c.b, func(t *testing.T) {
				tc.ShouldEvaluateTo(t, `return t.cmp("`+c.a+`",t:new("`+c.b+`"))`, c.expected)
			})
		}
	})

	t.Run("by_both_nil", func(t *testing.T) {
		tc.ShouldEvaluateTo(t, `return t.cmp(nil,nil)`, float64(0))
	})
	t.Run("by_n_nil", func(t *testing.T) {
		tc.ShouldEvaluateTo(t, `return t.cmp("0",nil)`, float64(1))
	})
	t.Run("by_a_nil", func(t *testing.T) {
		tc.ShouldEvaluateTo(t, `return t.cmp(nil,"0")`, float64(-1))
	})
}
