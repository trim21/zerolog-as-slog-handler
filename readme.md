# use zerolog as slog handler


```golang
package main

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"

	zerologassloghandler "github.com/trim21/zerolog-as-slog-handler"
)

func main() {
	h := zerologassloghandler.FromZerolog(zerolog.New(os.Stderr))

	l := slog.New(h)

	// {"level":"info","hello":"world","time":"2024-08-13T04:33:58+08:00","message":"1"}
	l.With("hello", "world").Info("1")

	// {"level":"info","j":{"hello":"world"},"time":"2024-08-13T04:33:58+08:00","message":"2"}
	l.WithGroup("j").With("hello", "world").Info("2")

	// {"level":"info","time":"2024-08-13T04:33:58+08:00","message":"3"}
	l.Info("3")
}
```
