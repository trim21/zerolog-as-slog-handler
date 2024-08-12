package zerologassloghandler

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/rs/zerolog"
)

type handler struct {
	groups []string
	log    zerolog.Logger
}

func (h handler) Enabled(ctx context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelDebug:
		return h.log.Debug().Enabled()
	case slog.LevelInfo:
		return h.log.Info().Enabled()
	case slog.LevelWarn:
		return h.log.Warn().Enabled()
	case slog.LevelError:
		return h.log.Error().Enabled()
	}

	return false
}

func (h handler) Handle(ctx context.Context, record slog.Record) error {
	var msg *zerolog.Event
	switch record.Level {
	case slog.LevelDebug:
		msg = h.log.Debug()
	case slog.LevelInfo:
		msg = h.log.Info()
	case slog.LevelWarn:
		msg = h.log.Warn()
	case slog.LevelError:
		msg = h.log.Error()
	default:
		panic(fmt.Sprint("unknown event"))
	}

	if msg == nil {
		return nil
	}

	if record.Time.IsZero() {
		msg = msg.Timestamp()
	} else {
		msg.Time(zerolog.TimestampFieldName, record.Time)
	}

	record.Attrs(func(attr slog.Attr) bool {
		msg = addToEvent(msg, attr.Key, attr.Value)
		return true
	})

	msg.Msg(record.Message)

	return nil
}

func addToEvent(ctx *zerolog.Event, key string, value slog.Value) *zerolog.Event {
	switch value.Kind() {
	case slog.KindAny:
		ctx = ctx.Any(key, value.Any())
	case slog.KindBool:
		ctx = ctx.Bool(key, value.Bool())
	case slog.KindDuration:
		ctx = ctx.Dur(key, value.Duration())
	case slog.KindFloat64:
		ctx = ctx.Float64(key, value.Float64())
	case slog.KindInt64:
		ctx = ctx.Int64(key, value.Int64())
	case slog.KindString:
		ctx = ctx.Str(key, value.String())
	case slog.KindTime:
		ctx = ctx.Time(key, value.Time())
	case slog.KindUint64:
		ctx = ctx.Uint64(key, value.Uint64())
	case slog.KindGroup:
		ctx = ctx.Any(key, value.Any()) // TODO
	case slog.KindLogValuer:
		return addToEvent(ctx, key, value.LogValuer().LogValue())
	default:
		ctx = ctx.Any(key, value.Any())
	}
	return ctx
}

func addToContext(ctx zerolog.Context, key string, value slog.Value) zerolog.Context {
	switch value.Kind() {
	case slog.KindAny:
		ctx = ctx.Any(key, value.Any())
	case slog.KindBool:
		ctx = ctx.Bool(key, value.Bool())
	case slog.KindDuration:
		ctx = ctx.Dur(key, value.Duration())
	case slog.KindFloat64:
		ctx = ctx.Float64(key, value.Float64())
	case slog.KindInt64:
		ctx = ctx.Int64(key, value.Int64())
	case slog.KindString:
		ctx = ctx.Str(key, value.String())
	case slog.KindTime:
		ctx = ctx.Time(key, value.Time())
	case slog.KindUint64:
		ctx = ctx.Uint64(key, value.Uint64())
	case slog.KindGroup:
		ctx = ctx.Any(key, value.Any()) // TODO
	case slog.KindLogValuer:
		return addToContext(ctx, key, value.LogValuer().LogValue())
	default:
		ctx = ctx.Any(key, value.Any())
	}
	return ctx
}

func (h handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return handler{
		groups: append(slices.Clone(h.groups), name), // TODO: perf
		log:    h.log,
	}
}

func (h handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(h.groups) == 0 {
		ctx := h.log.With()
		for _, attr := range attrs {
			ctx = addToContext(ctx, attr.Key, attr.Value)
		}
		h.log = ctx.Logger()
		return h
	}

	g := h.groups[len(h.groups)-1]

	return handler{
		groups: h.groups[:len(h.groups)-1],
		log:    h.log.With().Object(g, attrsObject{attrs: attrs}).Logger(),
	}
}

type attrsObject struct {
	attrs []slog.Attr
}

func (a attrsObject) MarshalZerologObject(e *zerolog.Event) {
	for _, i := range a.attrs {
		addToEvent(e, i.Key, i.Value)
	}
}

// FromZerolog create a slog.Handler from zerolog.Logger
// the logger **MUST** have timestamp disabled.
// for example
//
//	h := FromZerolog(zerolog.New(os.Stderr))
func FromZerolog(log zerolog.Logger) slog.Handler {
	return handler{log: log}
}
