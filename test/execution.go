package test

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/stretchr/testify/require"
)

func ExecMise(t testing.TB, args ...string) (string, int) {
	t.Helper()
	HookLogger(t)

	path, err := exec.LookPath("mise")
	require.NoError(t, err, "Should be able to find mise executable in PATH.")

	var output bytes.Buffer

	cmd := exec.Command(path, args...)
	cmd.Stdout = &output
	cmd.Stderr = &output

	err = cmd.Run()
	var exitErr *exec.ExitError
	var exitCode int
	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	} else {
		require.NoError(t, err, "Should be able to run mise without error.")
	}

	result := strings.TrimSpace(output.String())
	if result != "" {
		log.With("cmd", append([]string{"mise"}, args...)).
			With("exitCode", exitCode).
			Trace(result)
	}
	return result, exitCode
}

func ShouldExecMise(t testing.TB, expectedCode int, args ...string) string {
	t.Helper()
	output, code := ExecMise(t, args...)
	if code != expectedCode {
		t.Errorf("stdout/stderr:\n%s", output)
		require.Equal(t, expectedCode, code, "mise %v should exit with %d", args, expectedCode)
	}
	return output
}

func ShouldExecMiseMatching(t testing.TB, expectedCode int, contentMatching string, args ...string) {
	t.Helper()
	output := ShouldExecMise(t, expectedCode, args...)
	require.Regexp(t, contentMatching, output, "mise %v should match %s", args, contentMatching)
}
