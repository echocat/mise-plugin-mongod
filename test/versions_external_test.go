//go:build test_external

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersions_fetch_windows_External(t *testing.T) {
	tc := GivenContextWith(t, "../lib/versions.lua")
	tc.OsType = "windows"
	tc.ArchType = "amd64"

	actual := tc.ShouldEvaluate(t, `return t.__fetch()`)
	require.IsType(t, map[string]any{}, actual)

	actualVersionPlain := actual.(map[string]any)["8.0.0"]
	require.IsType(t, map[string]any{}, actualVersionPlain)
	actualVersion := actualVersionPlain.(map[string]any)

	assert.Equal(t, "lts", actualVersion["note"])
	assert.Equal(t, "https://docs.mongodb.org/master/release-notes/8.0/", actualVersion["release_notes"])
	assert.Equal(t, "base", actualVersion["edition"])
	assert.Equal(t, "https://fastdl.mongodb.org/windows/mongodb-windows-x86_64-8.0.0.zip", actualVersion["url"])
	assert.Equal(t, "8f7c86737cda331c5ca9491c64707d887d69cb3b", actualVersion["sha1"])
	assert.Equal(t, "4745e9d31b9414a0c708630768532797578df705107604c69b27ebb679c4b595", actualVersion["sha256"])
}

func TestVersions_fetch_macos_External(t *testing.T) {
	tc := GivenContextWith(t, "../lib/versions.lua")
	tc.OsType = "darwin"
	tc.ArchType = "arm64"

	actual := tc.ShouldEvaluate(t, `return t.__fetch()`)
	require.IsType(t, map[string]any{}, actual)

	actualVersionPlain := actual.(map[string]any)["8.0.0"]
	require.IsType(t, map[string]any{}, actualVersionPlain)
	actualVersion := actualVersionPlain.(map[string]any)

	assert.Equal(t, "lts", actualVersion["note"])
	assert.Equal(t, "https://docs.mongodb.org/master/release-notes/8.0/", actualVersion["release_notes"])
	assert.Equal(t, "base", actualVersion["edition"])
	assert.Equal(t, "https://fastdl.mongodb.org/osx/mongodb-macos-arm64-8.0.0.tgz", actualVersion["url"])
	assert.Equal(t, "e5ec7dc819d492b4dd3ae8392b7c0248443822e7", actualVersion["sha1"])
	assert.Equal(t, "4e51865ebe360b166045028622e49952412254548e0cc8825c3b84145717861c", actualVersion["sha256"])
}

func TestVersions_fetch_ubuntu2402_External(t *testing.T) {
	tc := GivenContextWith(t, "../lib/versions.lua")
	tc.OsType = "linux"
	tc.ArchType = "arm64"
	tc.DistributionType = "ubuntu"
	tc.DistributionVersion = "24.4"

	actual := tc.ShouldEvaluate(t, `return t.__fetch()`)
	require.IsType(t, map[string]any{}, actual)

	actualVersionPlain := actual.(map[string]any)["8.0.0"]
	require.IsType(t, map[string]any{}, actualVersionPlain)
	actualVersion := actualVersionPlain.(map[string]any)

	assert.Equal(t, "lts", actualVersion["note"])
	assert.Equal(t, "https://docs.mongodb.org/master/release-notes/8.0/", actualVersion["release_notes"])
	assert.Equal(t, "targeted", actualVersion["edition"])
	assert.Equal(t, "https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-ubuntu2004-8.0.0.tgz", actualVersion["url"])
	assert.Equal(t, "45d10603c349538b27d5399317ec1d29af92bdb8", actualVersion["sha1"])
	assert.Equal(t, "2ba041a9c5de5271b65969e201e2dcdb773fa45d0de1e1d52cf937592d315a73", actualVersion["sha256"])
}

func TestVersions_fetch_debian12_External(t *testing.T) {
	tc := GivenContextWith(t, "../lib/versions.lua")
	tc.OsType = "linux"
	tc.ArchType = "amd64"
	tc.DistributionType = "debian"
	tc.DistributionVersion = "12"

	actual := tc.ShouldEvaluate(t, `return t.__fetch()`)
	require.IsType(t, map[string]any{}, actual)

	actualVersionPlain := actual.(map[string]any)["8.0.0"]
	require.IsType(t, map[string]any{}, actualVersionPlain)
	actualVersion := actualVersionPlain.(map[string]any)

	assert.Equal(t, "lts", actualVersion["note"])
	assert.Equal(t, "https://docs.mongodb.org/master/release-notes/8.0/", actualVersion["release_notes"])
	assert.Equal(t, "targeted", actualVersion["edition"])
	assert.Equal(t, "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian12-8.0.0.tgz", actualVersion["url"])
	assert.Equal(t, "2ebc454354430dd9b73c931111d74f44e50871a9", actualVersion["sha1"])
	assert.Equal(t, "1743686860595bd194a60a5852d1c9447c3496581ede75e5ac495c22a481408e", actualVersion["sha256"])
}

func TestVersions_fetch_debian13_External(t *testing.T) {
	tc := GivenContextWith(t, "../lib/versions.lua")
	tc.OsType = "linux"
	tc.ArchType = "amd64"
	tc.DistributionType = "debian"
	tc.DistributionVersion = "13"

	actual := tc.ShouldEvaluate(t, `return t.__fetch()`)
	require.IsType(t, map[string]any{}, actual)

	actualVersionPlain := actual.(map[string]any)["8.0.0"]
	require.IsType(t, map[string]any{}, actualVersionPlain)
	actualVersion := actualVersionPlain.(map[string]any)

	assert.Equal(t, "lts", actualVersion["note"])
	assert.Equal(t, "https://docs.mongodb.org/master/release-notes/8.0/", actualVersion["release_notes"])
	assert.Equal(t, "targeted", actualVersion["edition"])
	assert.Equal(t, "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian12-8.0.0.tgz", actualVersion["url"])
	assert.Equal(t, "2ebc454354430dd9b73c931111d74f44e50871a9", actualVersion["sha1"])
	assert.Equal(t, "1743686860595bd194a60a5852d1c9447c3496581ede75e5ac495c22a481408e", actualVersion["sha256"])
}
