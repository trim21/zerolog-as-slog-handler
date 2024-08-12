package main

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"

	zerologassloghandler "zerolog-as-slog-handler"
)

func main() {
	h := zerologassloghandler.FromZerolog(zerolog.New(os.Stderr))

	l := slog.New(h)

	l.With("hello", "world").Info("1")

	l.WithGroup("j").With("hello", "world").Info("2")

	l.Info("3")
}
