//go:build test_e2e

package test

import (
	"testing"
)

func TestE2E_install(t *testing.T) {
	t.Cleanup(func() {
		_, _ = ExecMise(t, "plugin", "uninstall", "mongod")
	})

	ShouldExecMiseMatching(t, 0, `^$`, "plugin", "link", "--force", "mongod", "..")
	ShouldExecMise(t, 0, "cache", "clear")
	ShouldExecMise(t, 0, "install", "mongod@latest")
	ShouldExecMiseMatching(t, 0, `db version v\d+\.\d+\.\d+`, "exec", "mongod@latest", "--", "mongod", "--version")
}
