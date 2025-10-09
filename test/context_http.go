package test

import (
	"io"
	"net/http"
	"time"

	log "github.com/echocat/slf4g"
	lua "github.com/yuin/gopher-lua"
)

type contextHttp struct{}

func (m *contextHttp) get(L *lua.LState) int {
	param := L.CheckTable(1)
	urlStr := param.RawGetString("url")
	if urlStr == lua.LNil {
		L.Push(lua.LNil)
		L.Push(lua.LString("url is required"))
		return 2
	}

	req, err := http.NewRequest("GET", urlStr.String(), nil)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	headersTable := param.RawGetString("headers")
	if headersTable != lua.LNil {
		if table, ok := headersTable.(*lua.LTable); ok {
			table.ForEach(func(key lua.LValue, value lua.LValue) {
				req.Header.Add(key.String(), value.String())
			})
		}
	}

	start := time.Now()
	logger := log.With("url", urlStr).
		With("method", "GET")

	logger.Debug("Executing HTTP request...")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	headers := L.NewTable()
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers.RawSetString(k, lua.LString(v[0]))
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	L.SetField(result, "body", lua.LString(body))
	L.SetField(result, "status_code", lua.LNumber(resp.StatusCode))
	L.SetField(result, "headers", headers)
	L.SetField(result, "content_length", lua.LNumber(resp.ContentLength))
	L.Push(result)

	logger.
		With("duration", time.Since(start).Truncate(time.Millisecond).String()).
		Debug("Executing HTTP request... DONE!")

	return 1
}

func createContextHttpLoader() lua.LGFunction {
	return func(L *lua.LState) int {
		m := &contextHttp{}
		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{
			"get": m.get,
		})
		L.Push(t)
		return 1
	}
}
