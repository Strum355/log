package log

import (
	"fmt"
	"sort"
	"strings"
)

type Fields map[string]interface{}

func WithFields(f Fields) *Entry {
	return &Entry{fields: f, callDepth: callDepth - 1}
}

func (e *Entry) WithFields(f Fields) *Entry {
	if e.fields == nil {
		e.fields = make(Fields)
	}
	for k, v := range f {
		e.fields[k] = v
	}
	return e
}

func (f Fields) format() string {
	if f == nil || len(f) == 0 {
		return ""
	}

	keys := make([]string, 0, len(f))
	for k := range f {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteRune('\t')

	var iterCount int
	for _, k := range keys {
		v := f[k]
		s := strings.TrimSpace(fmt.Sprintf("%s=%v", k, v))
		if iterCount < len(f)-1 {
			s += " "
		}
		builder.WriteString(s)
		iterCount++
	}
	builder.WriteRune('\n')
	return builder.String()
}
