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

func Exec(t testing.TB, prg string, args ...string) (string, int) {
	t.Helper()
	HookLogger(t)

	path, err := exec.LookPath(prg)
	require.NoError(t, err, "Should be able to find %s executable in PATH.", prg)

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
		require.NoError(t, err, "Should be able to run %s without error.", prg)
	}

	result := strings.TrimSpace(output.String())
	if result != "" {
		log.With("cmd", append([]string{"mise"}, args...)).
			With("exitCode", exitCode).
			Trace(result)
	}
	return result, exitCode
}

func ShouldExec(t testing.TB, expectedCode int, prg string, args ...string) string {
	t.Helper()
	output, code := Exec(t, prg, args...)
	if code != expectedCode {
		t.Errorf("stdout/stderr:\n%s", output)
		require.Equal(t, expectedCode, code, "%s %v should exit with %d", prg, args, expectedCode)
	}
	return output
}

func ShouldMatching(t testing.TB, expectedCode int, contentMatching string, prg string, args ...string) {
	t.Helper()
	output := ShouldExec(t, expectedCode, prg, args...)
	require.Regexp(t, contentMatching, output, "%s %v should match %s", args, prg, contentMatching)
}
