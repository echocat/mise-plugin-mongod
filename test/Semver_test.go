package test

import (
	"testing"
)

func TestSemver_new(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Semver.lua")

	cases := []struct {
		input       string
		expected    any
		expectedErr string
	}{
		{`"1"`, nil, ""},
		{`"1.2"`, nil, ""},
		{`"1.2.3"`, map[string]any{"major": float64(1), "minor": float64(2), "patch": float64(3)}, ""},
		{`"1.2.0"`, map[string]any{"major": float64(1), "minor": float64(2), "patch": float64(0)}, ""},
		{`"0.0.0"`, map[string]any{"major": float64(0), "minor": float64(0), "patch": float64(0)}, ""},
		{`1`, nil, "requires a string to create a semver from; but got number"},
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

func TestSemver_cmp(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Semver.lua")

	cases := []struct {
		a, b     string
		expected float64
	}{
		{"1", "1", 0},     // Because both invalid
		{"2", "1", 0},     // Because both invalid
		{"1.1", "1.1", 0}, // Because both invalid
		{"1.2", "1.0", 0}, // Because both invalid

		{"1.2.3", "1.2.3", 0},
		{"2.2.3", "1.2.3", 1},
		{"1.2.3", "2.2.3", -1},

		{"1.2.3", "1.2.3", 0},
		{"1.3.3", "1.2.3", 1},
		{"1.2.3", "1.3.3", -1},

		{"1.2.3", "1.2.3", 0},
		{"1.2.4", "1.2.3", 1},
		{"1.2.3", "1.2.4", -1},
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
				tc.ShouldEvaluateTo(t, `local instance = t:new("`+c.a+`")
if instance == nil then
	return 0
end
return instance:cmp("`+c.b+`")`, c.expected)
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
		tc.ShouldEvaluateTo(t, `return t.cmp("1.2.3",nil)`, float64(1))
	})
	t.Run("by_a_nil", func(t *testing.T) {
		tc.ShouldEvaluateTo(t, `return t.cmp(nil,"1.2.3")`, float64(-1))
	})
}
