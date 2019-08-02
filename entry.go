package log

import (
	"io"
	"os"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/bwmarrin/lit"
)

const (
	callDepth = 4
)

type Entry struct {
	callDepth int
	fields    Fields
	span      opentracing.Span
}

var emptyEntry = &Entry{callDepth: callDepth}

func WithError(err error) *Entry {
	return &Entry{
		fields: Fields{
			"error": err,
		},
		callDepth: callDepth - 1,
	}
}

func (e *Entry) WithError(err error) *Entry {
	return e.WithFields(Fields{
		"error": err,
	})
}

func WithSpan(span opentracing.Span) *Entry {
	return &Entry{
		span:      span,
		callDepth: callDepth - 1,
	}
}

func (e *Entry) WithSpan(span opentracing.Span) *Entry {
	e.span = span
	return e
}

func (e *Entry) Clone() *Entry {
	fields := make(Fields, len(e.fields))
	for k, v := range e.fields {
		fields[k] = v
	}
	return &Entry{
		fields:    fields,
		callDepth: e.callDepth,
		span:      e.span,
	}
}

// ResetCallDepth should be called if not chaining a log call after
// calling precursor methods
//
// eg:
// e := log.WithFields(...)
// e.ResetCallDepth()
// e.Info(...)
func (e *Entry) ResetCallDepth() {
	e.callDepth = callDepth
}

func Debug(format string, a ...interface{}) {
	emptyEntry.Debug(format, a...)
}

func (e *Entry) Debug(format string, a ...interface{}) {
	e.log(lit.LogDebug, format, a...)
}

func Info(format string, a ...interface{}) {
	emptyEntry.Info(format, a...)
}

func (e *Entry) Info(format string, a ...interface{}) {
	e.log(lit.LogInformational, format, a...)
}

func Warn(format string, a ...interface{}) {
	emptyEntry.Warn(format, a...)
}

func (e *Entry) Warn(format string, a ...interface{}) {
	e.log(lit.LogWarning, format, a...)
}

func Error(format string, a ...interface{}) {
	emptyEntry.Error(format, a...)
}

func (e *Entry) Error(format string, a ...interface{}) {
	e.log(lit.LogError, format, a...)
}

func (e *Entry) log(level int, format string, a ...interface{}) {
	builder := new(strings.Builder)
	lit.Custom(builder, level, e.callDepth, format, a...)
	builder.WriteString(e.fields.format())
	if level == lit.LogError && config.UseStdErr {
		io.WriteString(os.Stderr, builder.String())
	} else {
		io.WriteString(config.Output, builder.String())
	}
}
