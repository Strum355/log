package log_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Strum355/log"
)

var (
	logRegex = regexp.MustCompile(`^(?:\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} )\[(\w{4,5} ?)\] ([^:]+):\d+:(.+\(\)) (.+)$`)
)

func initWithBuilder(s *strings.Builder) {
	log.Init(log.Config{
		Output: s,
	})
}

func assureSingleNewline(s string, t *testing.T) bool {
	hasExtraNewline := strings.HasSuffix(strings.TrimSuffix(s, "\n"), "\n")
	if hasExtraNewline {
		t.Error("log without fields had multiple newlines")
	}
	return hasExtraNewline
}

func splitMessage(s string, t *testing.T) (level, file, function, message string) {
	s = strings.TrimSpace(s)
	matched := logRegex.FindStringSubmatch(s)
	if len(matched) < 5 {
		t.Fatalf("message '%s' didnt match regex", s)
	}
	level = matched[1]
	file = matched[2]
	function = matched[3]
	message = matched[4]
	return
}

func Test_NoFields(t *testing.T) {
	var b strings.Builder

	initWithBuilder(&b)

	t.Run("No Fields", func(t *testing.T) {
		tests := []struct {
			level    string
			file     string
			function string
			f        func(format string, a ...interface{})
		}{
			{
				level:    "ERROR",
				file:     "entry_test.go",
				function: "1()",
				f:        log.Error,
			},
			{
				level:    "INFO ",
				file:     "entry_test.go",
				function: "1()",
				f:        log.Info,
			},
			{
				level:    "DEBUG",
				file:     "entry_test.go",
				function: "1()",
				f:        log.Debug,
			},
			{
				level:    "WARN ",
				file:     "entry_test.go",
				function: "1()",
				f:        log.Warn,
			},
		}

		for _, test := range tests {
			t.Run(test.level, func(t *testing.T) {
				test.f("there was an error")

				out := b.String()

				//t.Log(b.String())

				assureSingleNewline(out, t)

				level, file, function, _ := splitMessage(out, t)

				if level != test.level {
					t.Errorf("expected level: '%s'. actual level: '%s'\n", test.level, level)
				}

				if file != test.file {
					t.Errorf("expected file: '%s'. actual file: '%s'\n", test.file, file)
				}

				if function != test.function {
					t.Errorf("expected function: '%s'. actual function: '%s'\n", test.function, function)
				}

				b.Reset()
			})
		}
	})
}
