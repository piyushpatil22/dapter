package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

}

func Logger() zerolog.Logger{
	return Log
}