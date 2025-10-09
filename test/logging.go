package test

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/sdk/testlog"
)

func HookLogger(t testing.TB) {
	t.Helper()
	p := log.GetProvider()
	if w, ok := p.(interface{ Unwrap() log.Provider }); ok {
		if wp := w.Unwrap(); wp != nil {
			p = wp
		}
	}

	if _, alreadyHooked := p.(*testlog.Provider); alreadyHooked {
		return
	}

	provider := testlog.NewProvider(t)
	previous := log.SetProvider(provider)
	t.Cleanup(func() {
		log.SetProvider(previous)
	})
}
