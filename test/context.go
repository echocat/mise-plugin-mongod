package test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/echocat/slf4g"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

var (
	DefaultContextLogger = log.GetRootLogger()
	DefaultLibPath       = filepath.Join("..", "lib")

	stripErrorPrefixRegexp = regexp.MustCompile(`^[^:]+:\d+: `)
)

func GivenContext(t testing.TB) *Context {
	t.Helper()
	HookLogger(t)

	c := NewContext()
	t.Cleanup(c.Close)

	c.L.PreloadModule("http", createContextHttpLoader())
	c.L.PreloadModule("json", contextJsonLoader)

	if err := c.PreLoadLibDir(DefaultLibPath); err != nil {
		t.Fatal(err)
	}

	return c
}

func GivenContextWith(t testing.TB, luaFile string, otherModules ...string) *Context {
	t.Helper()
	c := GivenContext(t)
	L := c.getL()

	{
		lf, err := L.LoadFile(luaFile)
		if err != nil {
			t.Fatalf("cannot load lua file %q: %v", luaFile, err)
			return c
		}

		L.Push(lf)
		L.Call(0, 1)

		v := L.Get(-1)
		L.SetGlobal("t", v)
	}

	for _, om := range otherModules {
		lf, err := L.LoadFile(om)
		if err != nil {
			t.Fatalf("cannot load lua file %q: %v", om, err)
			return c
		}

		L.Push(lf)
		L.Call(0, 1)

		v := L.Get(-1)
		omn := strings.TrimSuffix(filepath.Base(om), filepath.Ext(om))
		L.SetGlobal(omn, v)
	}

	return c
}

func NewContext() *Context {
	L := lua.NewState()

	result := &Context{
		L:        L,
		OsType:   "Windows",
		ArchType: "amd64",
	}

	rt := L.NewTable()
	mt := L.NewTable()

	L.SetField(mt, "__index", L.NewFunction(func(L *lua.LState) int {
		// arg1 = das Table (RUNTIME), arg2 = key
		key := L.CheckString(2)
		switch key {
		case "osType":
			L.Push(lua.LString(result.OsType))
			return 1
		case "archType":
			L.Push(lua.LString(result.ArchType))
			return 1

		case "distributionType":
			if v := result.DistributionType; v != "" {
				L.Push(lua.LString(v))
				return 1
			}
		case "distributionVersion":
			if v := result.DistributionVersion; v != "" {
				L.Push(lua.LString(v))
				return 1
			}
		}
		L.Push(lua.LNil)
		return 1
	}))

	L.SetField(mt, "__newindex", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(2)
		L.RaiseError("RUNTIME is read-only (attempt to set %q)", key)
		return 0
	}))

	L.SetMetatable(rt, mt)
	L.SetGlobal("RUNTIME", rt)

	return result
}

type Context struct {
	L      *lua.LState
	Logger log.Logger

	OsType   string
	ArchType string

	DistributionType    string
	DistributionVersion string
}

func (c *Context) ShouldEvaluate(t testing.TB, source string) any {
	t.Helper()
	L := c.getL()
	fn, err := L.LoadString(source)
	require.NoError(t, err, "Evaluation of script should not fail.")

	L.Push(fn)
	err = L.PCall(0, 1, nil)
	require.NoError(t, err, "Execution of script should not fail.")

	lCurrent := L.Get(-1)
	result, err := c.ValueToAny(lCurrent)
	require.NoError(t, err, "Convertion of scripts execution result should not fail.")

	L.Pop(1)

	return result
}

func (c *Context) ShouldEvaluateTo(t testing.TB, source string, expected any) {
	t.Helper()
	actual := c.ShouldEvaluate(t, source)
	require.Equal(t, expected, actual, "Evaluation of %q should match %v", source, expected)
}

func (c *Context) ShouldEvaluateToError(t testing.TB, source string, expectedErrorContains string) {
	t.Helper()
	L := c.getL()
	fn, err := L.LoadString(source)
	require.NoError(t, err, "Evaluation of script should not fail.")

	L.Push(fn)
	err = L.PCall(0, 1, nil)

	tErr := err
	var lae *lua.ApiError
	if errors.As(tErr, &lae) && lae.Object != nil {
		tErr = errors.New(stripErrorPrefixRegexp.ReplaceAllString(lae.Object.String(), ""))
	}

	require.ErrorContains(t, tErr, expectedErrorContains, "Evaluation of %q should fail with an error containing %q", source, expectedErrorContains)
}

func (c *Context) PreLoadLibDir(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}

	L := c.getL()
	des, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("pre-load lib dir %q: %w", path, err)
	}
	for _, de := range des {
		if !de.Type().IsRegular() {
			continue
		}
		den := de.Name()
		def := filepath.Join(path, den)
		den = strings.TrimSuffix(den, filepath.Ext(den))
		L.PreloadModule(den, func(L *lua.LState) int {
			lf, err := L.LoadFile(def)
			logger := c.GetLogger().
				With("module", den)
			if err != nil {
				L.RaiseError("cannot load module %q: %v", den, err)
				logger.WithError(err).Trace("contextHttp cannot be loaded.")
				return 1
			}

			L.Push(lf)
			L.Call(0, 1)
			logger.Trace("contextHttp loaded.")
			return 1
		})
	}

	return nil
}

func (c *Context) GetLogger() log.Logger {
	if c != nil {
		if v := c.Logger; v != nil {
			return v
		}
	}
	return DefaultContextLogger
}

func (c *Context) getL() *lua.LState {
	if c == nil {
		panic("nil context")
	}
	if v := c.L; v != nil {
		return v
	}
	panic("context not initialized")
}

func (c *Context) Close() {
	if c == nil {
		return
	}
	if l := c.L; l != nil {
		l.Close()
	}
}
