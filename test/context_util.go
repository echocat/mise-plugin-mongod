package test

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

func (c *Context) ValueToAny(v lua.LValue) (any, error) {
	seen := map[*lua.LTable]struct{}{}
	return c.cvt(v, 0, seen)
}

func (c *Context) cvt(v lua.LValue, depth int, seen map[*lua.LTable]struct{}) (any, error) {
	if depth > 25 {
		return nil, fmt.Errorf("max depth reached: %d", depth)
	}

	switch v.Type() {
	case lua.LTNil:
		return nil, nil
	case lua.LTBool:
		return bool(v.(lua.LBool)), nil
	case lua.LTNumber:
		n := float64(v.(lua.LNumber))
		return n, nil
	case lua.LTString:
		return string(v.(lua.LString)), nil
	case lua.LTUserData:
		ud := v.(*lua.LUserData)
		if ud.Value != nil {
			return ud.Value, nil
		}
		return nil, nil
	case lua.LTTable:
		t := v.(*lua.LTable)
		if _, ok := seen[t]; ok {
			return nil, errors.New("cycle detected in table")
		}
		seen[t] = struct{}{}
		defer delete(seen, t)

		n := t.Len()
		if c.isPureArray(t, n) {
			out := make([]any, n)
			for i := 1; i <= n; i++ {
				vi, err := c.cvt(t.RawGetInt(i), depth+1, seen)
				if err != nil {
					return nil, err
				}
				out[i-1] = vi
			}
			return out, nil
		}

		// Map[string]any – Schlüssel werden zu Strings (Lua tostring, respektiert __tostring)
		out := make(map[string]any)
		t.ForEach(func(k, val lua.LValue) {
			// leer – wir sammeln über Closure
		})
		// Nochmal iterieren, weil wir Fehler behandeln wollen:
		var err error
		t.ForEach(func(k, val lua.LValue) {
			if err != nil {
				return
			}
			ks, kerr := c.luaValueToString(k)
			if kerr != nil {
				err = kerr
				return
			}
			gv, verr := c.cvt(val, depth+1, seen)
			if verr != nil {
				err = verr
				return
			}
			out[ks] = gv
		})
		if err != nil {
			return nil, err
		}
		return out, nil

	default:
		// Funktionen, Threads, Channels etc.
		return nil, errors.New("unsupported lua type: " + v.Type().String())
	}
}

func (c *Context) isPureArray(t *lua.LTable, n int) bool {
	count := 0
	pure := true
	t.ForEach(func(k, _ lua.LValue) {
		if !pure {
			return
		}
		if kn, ok := k.(lua.LNumber); ok {
			f := float64(kn)
			if f == math.Trunc(f) && int(f) >= 1 && int(f) <= n {
				count++
				return
			}
		}
		pure = false
	})
	return pure && count == n
}

func (c *Context) luaValueToString(v lua.LValue) (string, error) {
	switch x := v.(type) {
	case lua.LString:
		return string(x), nil
	case lua.LNumber:
		f := float64(x)
		if f == math.Trunc(f) {
			return strconv.FormatInt(int64(f), 10), nil
		}
		return strconv.FormatFloat(f, 'g', -1, 64), nil
	case lua.LBool:
		if bool(x) {
			return "true", nil
		}
		return "false", nil
	}

	L := c.getL()

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("tostring"),
		NRet:    1,
		Protect: true,
	}, v); err != nil {
		return "", err
	}
	s := L.ToString(-1)
	L.Pop(1)
	return s, nil
}
