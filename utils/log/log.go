package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"strconv"
	"time"
)

func init() {

	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
}

func Error() *zerolog.Event {

	_, file, line, ok := runtime.Caller(1)
	e := log.Error()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Debug() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Debug()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Warn() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Warn()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}

func Info() *zerolog.Event {
	_, file, line, ok := runtime.Caller(1)
	e := log.Info()
	if ok {
		e = e.Str("line", file+":"+strconv.Itoa(line))
	}
	return e
}
