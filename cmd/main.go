package main

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"

	zerologassloghandler "github.com/trim21/zerolog-as-slog-handler"
)

func main() {
	h := zerologassloghandler.FromZerolog(zerolog.New(os.Stderr))

	{
		l := slog.New(slog.NewJSONHandler(os.Stderr, nil))
		l.With("a", "b").WithGroup("G").With("c", "d").WithGroup("H").Info("msg", "e", "f")
	}

	{
		l := slog.New(h)
		l.With("a", "b").WithGroup("G").With("c", "d").WithGroup("H").Info("msg", "e", "f")
	}
}
