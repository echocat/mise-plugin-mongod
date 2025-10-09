package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTarget_new(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	t.Run("fromString", func(t *testing.T) {
		cases := []struct {
			input       string
			expected    any
			expectedErr string
		}{
			{"windows", map[string]any{"os": "windows"}, ""},
			{"windows13", nil, "Unknown target windows13."},
			{"windows2404", nil, "Unknown target windows2404."},

			{"macos", map[string]any{"os": "macos"}, ""},
			{"macos13", nil, "Unknown target macos13."},
			{"macos2404", nil, "Unknown target macos2404."},

			{"ubuntu2402", map[string]any{"os": "linux", "distribution": "ubuntu", "version": []any{float64(24), float64(2)}}, ""},
			{"ubuntu2202", map[string]any{"os": "linux", "distribution": "ubuntu", "version": []any{float64(22), float64(2)}}, ""},
			{"ubuntu0102", map[string]any{"os": "linux", "distribution": "ubuntu", "version": []any{float64(1), float64(2)}}, ""},
			{"ubuntu13", nil, "Version of target ubuntu13 cannot be interpreted."},

			{"debian13", map[string]any{"os": "linux", "distribution": "debian", "version": []any{float64(13)}}, ""},
			{"debian666", map[string]any{"os": "linux", "distribution": "debian", "version": []any{float64(666)}}, ""},
			{"debian2202", map[string]any{"os": "linux", "distribution": "debian", "version": []any{float64(2202)}}, ""},

			{"suse13", map[string]any{"os": "linux", "distribution": "suse", "version": []any{float64(13)}}, ""},
			{"suse666", map[string]any{"os": "linux", "distribution": "suse", "version": []any{float64(666)}}, ""},
			{"suse2202", map[string]any{"os": "linux", "distribution": "suse", "version": []any{float64(2202)}}, ""},

			{"rhel83", map[string]any{"os": "linux", "distribution": "rhel", "version": []any{float64(8), float64(3)}}, ""},
			{"rhel9", map[string]any{"os": "linux", "distribution": "rhel", "version": []any{float64(9)}}, ""},
			{"rhel123", nil, "Version of target rhel123 cannot be interpreted."},
		}

		for _, c := range cases {
			t.Run(c.input, func(t *testing.T) {
				if expectedErr := c.expectedErr; expectedErr == "" {
					tc.ShouldEvaluateTo(t, `return t:new("`+c.input+`")`, c.expected)
				} else {
					tc.ShouldEvaluateToError(t, `return t:new("`+c.input+`")`, expectedErr)
				}
			})
		}
	})
}

func TestTarget_equals_base(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		a, b     string
		expected bool
	}{
		{"ubuntu2402", "ubuntu2402", true},
		{"ubuntu2402", "ubuntu2202", true},
		{"ubuntu2202", "ubuntu2402", true},
		{"ubuntu2202", "debian13", false},
		{"debian13", "ubuntu2202", false},
	}

	for _, c := range cases {
		t.Run(c.a+"_"+c.b, func(t *testing.T) {
			tc.ShouldEvaluateTo(t, `return t:new("`+c.a+`"):equals_base(t:new("`+c.b+`"))`, c.expected)
		})
	}
}

func TestTarget___parse_version_year_month(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		input    string
		expected any
	}{
		{`"2402"`, []any{float64(24), float64(2)}},
		{`"0102"`, []any{float64(1), float64(2)}},
		{`"666"`, nil},
		{`nil`, nil},
		{`123`, nil},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tc.ShouldEvaluateTo(t, `return t.__parse_version_year_month(`+c.input+`)`, c.expected)
		})
	}
}

func TestTarget___format_version_year_month(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		input       string
		expected    any
		expectedErr string
	}{
		{`{24,2}`, "2402", ``},
		{`{1,2}`, "0102", ``},
		{`nil`, "", ``},
		{`{}`, "", ``},
		{`{1}`, "", `Should format a version year month; but got: table: 0x`},
		{`{1,-1}`, "", `Should format a version year month; but got: table: 0x`},
		{`{-1,1}`, "", `Should format a version year month; but got: table: 0x`},
		{`{-1,-1}`, "", `Should format a version year month; but got: table: 0x`},
		{`{1,2,3}`, "", `Should format a version year month; but got: table: 0x`},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if expectedErr := c.expectedErr; expectedErr == "" {
				tc.ShouldEvaluateTo(t, `return t.__format_version_year_month(`+c.input+`)`, c.expected)
			} else {
				tc.ShouldEvaluateToError(t, `return t.__format_version_year_month(`+c.input+`)`, expectedErr)
			}
		})
	}
}

func TestTarget___parse_version_major_only(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		input    string
		expected any
	}{
		{`"1"`, []any{float64(1)}},
		{`"666"`, []any{float64(666)}},
		{`"666a"`, nil},
		{`nil`, nil},
		{`123`, nil},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tc.ShouldEvaluateTo(t, `return t.__parse_version_major_only(`+c.input+`)`, c.expected)
		})
	}
}

func TestTarget___format_version_major_only(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		input       string
		expected    any
		expectedErr string
	}{
		{`{1}`, "1", ``},
		{`{666}`, "666", ``},
		{`nil`, "", ``},
		{`{}`, "", ``},
		{`{-1}`, "", `Should format a major only version; but got: table: 0x`},
		{`{1,2}`, "", `Should format a major only version; but got: table: 0x`},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if expectedErr := c.expectedErr; expectedErr == "" {
				tc.ShouldEvaluateTo(t, `return t.__format_version_major_only(`+c.input+`)`, c.expected)
			} else {
				tc.ShouldEvaluateToError(t, `return t.__format_version_major_only(`+c.input+`)`, expectedErr)
			}
		})
	}
}

func TestTarget___parse_version_rhel(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		input    string
		expected any
	}{
		{`"83"`, []any{float64(8), float64(3)}},
		{`"9"`, []any{float64(9)}},
		{`"123"`, nil},
		{`"666a"`, nil},
		{`nil`, nil},
		{`123`, nil},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tc.ShouldEvaluateTo(t, `return t.__parse_version_rhel(`+c.input+`)`, c.expected)
		})
	}
}

func TestTarget___format_version_rhel(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	cases := []struct {
		input       string
		expected    any
		expectedErr string
	}{
		{`{8,3}`, "83", ``},
		{`{8,0}`, "8", ``},
		{`{9}`, "9", ``},
		{`nil`, "", ``},
		{`{}`, "", ``},
		{`{-1}`, "", `Should format a rhel version; but got: table: 0x`},
		{`{10}`, "", `Should format a rhel version; but got: table: 0x`},
		{`{8,-1}`, "", `Should format a rhel version; but got: table: 0x`},
		{`{8,10}`, "810", ``},
		{`{1,2,3}`, "", `Should format a rhel version; but got: table: 0x`},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if expectedErr := c.expectedErr; expectedErr == "" {
				tc.ShouldEvaluateTo(t, `return t.__format_version_rhel(`+c.input+`)`, c.expected)
			} else {
				tc.ShouldEvaluateToError(t, `return t.__format_version_rhel(`+c.input+`)`, expectedErr)
			}
		})
	}
}

func TestTarget__tostring(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua", "../lib/Version.lua")

	cases := []struct {
		input       string
		expected    string
		expectedErr string
	}{
		{`{os = "windows"}`, "windows", ""},
		{`{os = "macos"}`, "macos", ""},
		{`{os = "foobar"}`, "foobar", ""},
		{`{os = "linux", distribution = "ubuntu"}`, "ubuntu", ""},
		{`{os = "linux", distribution = "ubuntu", version = Version:new("22.2")}`, "ubuntu2202", ""},
		{`{os = "linux", distribution = "debian", version = Version:new("13")}`, "debian13", ""},
		{`{os = "linux", distribution = "suse", version = Version:new("13")}`, "suse13", ""},
		{`{os = "linux", distribution = "foo", version = Version:new("12.34")}`, "foo12.34", ""},
		{`{os = "linux", distribution = "rhel", version = Version:new("8.3")}`, "rhel83", ""},
		{`{os = "linux", distribution = "rhel", version = Version:new("9")}`, "rhel9", ""},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if expectedErr := c.expectedErr; expectedErr == "" {
				tc.ShouldEvaluateTo(t, `return tostring(t:new(`+c.input+`))`, c.expected)
			} else {
				tc.ShouldEvaluateToError(t, `return tostring(t:new(`+c.input+`))`, expectedErr)
			}
		})
	}

}

func TestTarget_instanceof(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua", "../lib/types.lua")

	cases := []struct {
		givenObj   string
		givenClass string
		expected   bool
	}{
		{`t:new("ubuntu2202")`, `t`, true},
		{`t:new("ubuntu2202")`, `nil`, false},
	}

	for _, c := range cases {
		t.Run(c.givenObj+"_"+c.givenClass, func(t *testing.T) {
			tc.ShouldEvaluateTo(t, `return types.instanceof(`+c.givenObj+`, `+c.givenClass+`)`, c.expected)
		})
	}
}

func TestTarget_host(t *testing.T) {
	tc := GivenContextWith(t, "../lib/Target.lua")

	t.Run("fromResolution", func(t *testing.T) {
		cases := []struct {
			name           string
			givenOs        string
			givenOsRelease string
			expected       any
			expectedErr    string
		}{
			{"windows", "windows", "", map[string]any{"os": "windows"}, ""},
			{"macos", "macos", "", map[string]any{"os": "macos"}, ""},
			{"wrong-os", "does-not-exist", "", nil, `Unsupported operating system: does-not-exist`},
			{"ubuntu2402", "linux", etcOsReleaseUbuntu2402, map[string]any{"distribution": "ubuntu", "os": "linux", "version": []any{float64(24), float64(4)}}, ""},
			{"debian13", "linux", etcOsReleaseDebian13, map[string]any{"distribution": "debian", "os": "linux", "version": []any{float64(13)}}, ""},
			{"amazon2023", "linux", etcOsReleaseAmazon2023, map[string]any{"distribution": "amazon", "os": "linux", "version": []any{float64(2023)}}, ""},
			{"rhel8", "linux", etcOsReleaseRhel8, map[string]any{"distribution": "rhel", "os": "linux", "version": []any{float64(8), float64(10)}}, ""},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				var osReleaseFn string

				if v := c.givenOsRelease; v != "" {
					osReleaseFn = filepath.Join(t.TempDir(), "os-release")
					require.NoError(t, os.WriteFile(osReleaseFn, []byte(v), 0600))
				}

				if expectedErr := c.expectedErr; expectedErr == "" {
					tc.ShouldEvaluateTo(t, `return t.host("`+c.givenOs+`", [[`+osReleaseFn+`]])`, c.expected)
				} else {
					tc.ShouldEvaluateToError(t, `return t.host("`+c.givenOs+`", [[`+osReleaseFn+`]])`, expectedErr)
				}
			})
		}
	})

	t.Run("fromEnv", func(t *testing.T) {
		cases := []struct {
			given       string
			expected    any
			expectedErr string
		}{
			{"windows", map[string]any{"os": "windows"}, ""},
			{"macos", map[string]any{"os": "macos"}, ""},
			{"wrong-os", nil, `Unknown target wrong-os.`},
			{"ubuntu2402", map[string]any{"distribution": "ubuntu", "os": "linux", "version": []any{float64(24), float64(2)}}, ""},
			{"debian13", map[string]any{"distribution": "debian", "os": "linux", "version": []any{float64(13)}}, ""},
			{"amazon2023", map[string]any{"distribution": "amazon", "os": "linux", "version": []any{float64(2023)}}, ""},
			{"rhel8", map[string]any{"distribution": "rhel", "os": "linux", "version": []any{float64(8)}}, ""},
		}

		for _, c := range cases {
			t.Run(c.given, func(t *testing.T) {
				require.NoError(t, os.Setenv("MONGOD_TARGET", c.given))
				t.Cleanup(func() {
					_ = os.Unsetenv("MONGOD_TARGET")
				})

				if expectedErr := c.expectedErr; expectedErr == "" {
					tc.ShouldEvaluateTo(t, `return t.host()`, c.expected)
				} else {
					tc.ShouldEvaluateToError(t, `return t.host()`, expectedErr)
				}
			})
		}
	})
}

const (
	etcOsReleaseUbuntu2402 = `PRETTY_NAME="Ubuntu 24.04.2 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04.2 LTS (Noble Numbat)"
VERSION_CODENAME=noble
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=noble
LOGO=ubuntu-logo`
	etcOsReleaseDebian13 = `PRETTY_NAME="Debian GNU/Linux 13 (trixie)"
NAME="Debian GNU/Linux"
VERSION_ID="13"
VERSION="13 (trixie)"
VERSION_CODENAME=trixie
DEBIAN_VERSION_FULL=13.1
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/`
	etcOsReleaseAmazon2023 = `NAME="Amazon Linux"
VERSION="2023"
ID="amzn"
ID_LIKE="fedora"
VERSION_ID="2023"
PLATFORM_ID="platform:al2023"
PRETTY_NAME="Amazon Linux 2023.9.20250929"
ANSI_COLOR="0;33"
CPE_NAME="cpe:2.3:o:amazon:amazon_linux:2023"
HOME_URL="https://aws.amazon.com/linux/amazon-linux-2023/"
DOCUMENTATION_URL="https://docs.aws.amazon.com/linux/"
SUPPORT_URL="https://aws.amazon.com/premiumsupport/"
BUG_REPORT_URL="https://github.com/amazonlinux/amazon-linux-2023"
VENDOR_NAME="AWS"
VENDOR_URL="https://aws.amazon.com/"
SUPPORT_END="2029-06-30"`
	etcOsReleaseRhel8 = `NAME="Red Hat Enterprise Linux"
VERSION="8.10 (Ootpa)"
ID="rhel"
ID_LIKE="fedora"
VERSION_ID="8.10"
PLATFORM_ID="platform:el8"
PRETTY_NAME="Red Hat Enterprise Linux 8.10 (Ootpa)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:redhat:enterprise_linux:8::baseos"
HOME_URL="https://www.redhat.com/"
DOCUMENTATION_URL="https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8"
BUG_REPORT_URL="https://issues.redhat.com/"

REDHAT_BUGZILLA_PRODUCT="Red Hat Enterprise Linux 8"
REDHAT_BUGZILLA_PRODUCT_VERSION=8.10
REDHAT_SUPPORT_PRODUCT="Red Hat Enterprise Linux"
REDHAT_SUPPORT_PRODUCT_VERSION="8.10"`
)
