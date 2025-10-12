//go:build test_e2e

package test

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_mise_install(t *testing.T) {
	t.Cleanup(func() {
		Exec(t, "mise", "uninstall", "mongod")
		Exec(t, "mise", "plugin", "uninstall", "mongod")
	})

	ShouldExec(t, 0, "mise", "plugin", "link", "--force", "mongod", "..")
	ShouldExec(t, 0, "mise", "cache", "clear")
	ShouldExec(t, 0, "mise", "install", "mongod@latest")
	ShouldMatching(t, 0, `db version v\d+\.\d+\.\d+`, "mise", "exec", "mongod@latest", "--", "mongod", "--version")
}

func TestE2E_vfox_install(t *testing.T) {
	root, err := os.Getwd()
	require.NoError(t, err, "Should get current working directory.")
	root = filepath.Dir(root)

	uh, err := user.Current()
	require.NoError(t, err, "Should get current user home directory.")

	pluginDir := filepath.Join(uh.HomeDir, ".version-fox", "plugin")
	err = os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err, "Should create %s directory.", pluginDir)

	pluginMountDir := filepath.Join(pluginDir, "mongod")
	_ = os.Remove(pluginMountDir)
	err = os.Symlink(root, pluginMountDir)
	require.NoError(t, err, "Should create symlink %s => %s.", root, pluginMountDir)

	t.Cleanup(func() {
		Exec(t, "vfox", "uninstall", "mongod")
		_ = os.Remove(pluginMountDir)
	})

	ShouldExec(t, 0, "vfox", "install", "mongod@8.2.1")
	rspJson := ShouldExec(t, 0, "vfox", "env", "--json", "mongod@8.2.1")
	var rsp vfoxEnvResponse
	err = json.Unmarshal([]byte(rspJson), &rsp)
	require.NoError(t, err, "Should be able to unmarshal json.")
	require.Len(t, rsp.Paths, 1, "Should have the path for mongod, but was: %s", rspJson)
	mongodExe := filepath.Join(rsp.Paths[0], "mongod")

	ShouldMatching(t, 0, `db version v\d+\.\d+\.\d+`, mongodExe, "--version")
}

type vfoxEnvResponse struct {
	Paths []string `json:"paths"`
}
