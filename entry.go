package log

import (
	"context"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
)

type logKey struct{}

var Key logKey = struct{}{}

type Entry struct {
	fields Fields
	span   opentracing.Span
}

var emptyEntry = &Entry{}

func WithError(err error) *Entry {
	return &Entry{
		fields: Fields{
			"error": err,
		},
	}
}

func (e *Entry) WithError(err error) *Entry {
	return e.WithFields(Fields{
		"error": err,
	})
}

func (e *Entry) WithContext(ctx context.Context) *Entry {
	if fields, ok := ctx.Value(Key).(*Fields); ok {
		*e = *e.WithFields(*fields)
	}
	return e
}

func WithContext(ctx context.Context) *Entry {
	e := new(Entry)
	if fields, ok := ctx.Value(Key).(Fields); ok {
		e.fields = fields
	}
	return e
}

/* func WithSpan(span opentracing.Span) *Entry {
	return &Entry{
		span: span,
	}
}

func (e *Entry) WithSpan(span opentracing.Span) *Entry {
	e.span = span
	return e
} */

func (e *Entry) Clone() *Entry {
	fields := make(Fields, len(e.fields))
	for k, v := range e.fields {
		fields[k] = v
	}
	return &Entry{
		fields: fields,
		span:   e.span,
	}
}

func Debug(msg string) {
	emptyEntry.Debug(msg)
}

func (e *Entry) Debug(msg string) {
	e.log(LogDebug, msg)
}

func Info(msg string) {
	emptyEntry.Info(msg)
}

func (e *Entry) Info(msg string) {
	e.log(LogInformational, msg)
}

func Warn(msg string) {
	emptyEntry.Warn(msg)
}

func (e *Entry) Warn(msg string) {
	e.log(LogWarning, msg)
}

func Error(msg string) {
	emptyEntry.Error(msg)
}

func (e *Entry) Error(msg string) {
	e.log(LogError, msg)
}

func (e *Entry) log(level LogLevel, format string) {
	if level < config.LogLevel {
		return
	}

	now := time.Now()

	builder := new(strings.Builder)

	file, fileLine, funcName := getFunctionInfo()

	config.logger.createLogPoint(logPoint{builder, level, fileLine, file, funcName, format, e.fields, now})

	if level == LogError && config.UseStdErr {
		io.WriteString(os.Stderr, builder.String())
	} else {
		io.WriteString(config.Output, builder.String())
	}
}

func getPrefix(level LogLevel) string {
	if level == LogError {
		return config.ErrorPrefix
	} else if level == LogWarning {
		return config.WarnPrefix
	} else if level == LogInformational {
		return config.InfoPrefix
	}
	return config.DebugPrefix
}

func getFunctionInfo() (file string, line int, name string) {
	pc, _, _, _ := runtime.Caller(0)

	// get package name and filename
	ownPkgName := strings.Join(strings.Split(runtime.FuncForPC(pc).Name(), ".")[:2], ".")

	pkgName := ownPkgName

	callDepth := 1
	// find the first package that isnt this one, excluding our tests
	for pkgName == ownPkgName && !strings.HasSuffix(file, "_test.go") {
		pc, file, line, _ = runtime.Caller(callDepth)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		pkgName = strings.Join(strings.Split(runtime.FuncForPC(pc).Name(), ".")[:2], ".")

		name = runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]
		callDepth++
	}

	return
}
