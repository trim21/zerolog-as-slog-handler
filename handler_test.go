package zerologassloghandler

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/rs/zerolog"
)

func TestSlog(t *testing.T) {
	zerolog.MessageFieldName = slog.MessageKey
	zerolog.TimestampFieldName = slog.TimeKey

	var o = make(map[string]*bytes.Buffer)

	slogtest.Run(t, func(t *testing.T) slog.Handler {
		var w = bytes.NewBuffer(nil)
		o[t.Name()] = w
		return FromZerolog(zerolog.New(w))
	}, func(t *testing.T) map[string]any {
		t.Log(o[t.Name()].String())
		var r map[string]any
		err := json.NewDecoder(o[t.Name()]).Decode(&r)
		if err != nil {
			panic(err)
		}
		return r
	})
}
