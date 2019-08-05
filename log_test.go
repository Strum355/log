package log_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Strum355/log"
)

var (
	logRegex = regexp.MustCompile(`^(?:\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:Z|(?:\+|-)\d{2}:\d{2})) \[(\w{4,5} ?)\] ([^:]+):\d+:(.+\(\)) (.+)$`)
	b        = new(strings.Builder)
)

func assureSingleNewline(s string, t *testing.T) bool {
	hasExtraNewline := strings.HasSuffix(strings.TrimSuffix(s, "\n"), "\n")
	if hasExtraNewline {
		t.Error("log without fields had multiple newlines")
	}
	return hasExtraNewline
}

func splitMessage(s string, t *testing.T) (level, file, function, message string) {
	s = strings.TrimSpace(strings.Split(s, "\n")[0])
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

func hasField(k string, v interface{}, out string, t *testing.T) (bool, string) {
	out = strings.TrimSpace(strings.Split(out, "\n")[1])
	return strings.Contains(out, strings.TrimSpace(fmt.Sprintf("%s='%v'", k, v))), out
}

// Test_SimpleLogger tests the majority of the codepaths for entry.go, fields.go, log.go and simple.go
func Test_SimpleLogger(t *testing.T) {
	defer b.Reset()

	t.Run("NoFields", func(t *testing.T) {
		defer b.Reset()
		log.InitSimpleLogger(&log.Config{
			Output: b,
		})

		tests := []struct {
			level    string
			file     string
			function string
			f        func(msg string)
		}{
			{
				level:    "ERROR",
				file:     "log_test.go",
				function: "1()",
				f:        log.Error,
			},
			{
				level:    "INFO ",
				file:     "log_test.go",
				function: "1()",
				f:        log.Info,
			},
			{
				level:    "DEBUG",
				file:     "log_test.go",
				function: "1()",
				f:        log.Debug,
			},
			{
				level:    "WARN ",
				file:     "log_test.go",
				function: "1()",
				f:        log.Warn,
			},
		}

		for _, test := range tests {
			t.Run(test.level, func(t *testing.T) {
				test.f("there was an error")
				defer b.Reset()

				out := b.String()

				assureSingleNewline(out, t)

				level, file, function, _ := splitMessage(out, t)

				if level != test.level {
					t.Errorf("expected level: '%s'. actual level: '%s'", test.level, level)
				}

				if file != test.file {
					t.Errorf("expected file: '%s'. actual file: '%s'", test.file, file)
				}

				if function != test.function {
					t.Errorf("expected function: '%s'. actual function: '%s'", test.function, function)
				}

				if len(strings.Split(strings.TrimSpace(out), "\n")) > 1 {
					t.Errorf("expected single line log point: '%s", out)
				}
			})
		}
	})

	t.Run("WithFields", func(t *testing.T) {
		defer b.Reset()
		t.Run("Single Field", func(t *testing.T) {
			log.InitSimpleLogger(&log.Config{
				Output: b,
			})

			tests := []struct {
				level    string
				file     string
				function string
				key      string
				value    interface{}
				f        func(string)
			}{
				{
					level:    "ERROR",
					file:     "log_test.go",
					function: "1()",
					key:      "sample",
					value:    "banana",
					f:        log.WithFields(log.Fields{"sample": "banana"}).Error,
				},
				{
					level:    "INFO ",
					file:     "log_test.go",
					function: "1()",
					key:      "text",
					value:    1,
					f:        log.WithFields(log.Fields{"text": 1}).Info,
				},
				{
					level:    "DEBUG",
					file:     "log_test.go",
					function: "1()",
					key:      "burger",
					value:    []string{"sorry fellas"},
					f:        log.WithFields(log.Fields{"burger": []string{"sorry fellas"}}).Debug,
				},
				{
					level:    "WARN ",
					file:     "log_test.go",
					function: "1()",
					key:      "salad",
					value:    "fortnite",
					f:        log.WithFields(log.Fields{"salad": "fortnite"}).Warn,
				},
			}

			for _, test := range tests {
				t.Run(test.level, func(t *testing.T) {
					test.f("there was an error")
					defer b.Reset()

					out := b.String()

					//t.Log(b.String())

					assureSingleNewline(out, t)

					level, file, function, _ := splitMessage(out, t)

					if level != test.level {
						t.Errorf("expected level: '%s'. actual level: '%s'", test.level, level)
					}

					if file != test.file {
						t.Errorf("expected file: '%s'. actual file: '%s'", test.file, file)
					}

					if function != test.function {
						t.Errorf("expected function: '%s'. actual function: '%s'", test.function, function)
					}

					if ok, fields := hasField(test.key, test.value, out, t); !ok {
						t.Errorf("expected fields to contain: '%s=%v. actual fields total: %s", test.key, test.value, fields)
					}
				})
			}
		})

		t.Run("Multiple Fields", func(t *testing.T) {
			defer b.Reset()
			log.InitSimpleLogger(&log.Config{
				Output: b,
			})

			tests := []struct {
				level    string
				file     string
				function string
				fields   log.Fields
				f        func(string)
			}{
				{
					level:    "ERROR",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"one":   1,
						"two":   "2",
						"three": []string{"1", "2", "3"},
					},
					f: log.WithFields(log.Fields{
						"one":   1,
						"two":   "2",
						"three": []string{"1", "2", "3"},
					}).Error,
				},
				{
					level:    "INFO ",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"sample": "this is a long piece of text",
						"true":   false,
						"false":  true,
					},
					f: log.WithFields(log.Fields{
						"sample": "this is a long piece of text",
						"true":   false,
						"false":  true,
					}).Info,
				},
				{
					level:    "WARN ",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"one":      nil,
						"okay but": "epic",
					},
					f: log.WithFields(log.Fields{
						"one":      nil,
						"okay but": "epic",
					}).Warn,
				},
				{
					level:    "DEBUG",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"teamwork":  -1,
						"dreamwork": []bool{false, true},
					},
					f: log.WithFields(log.Fields{
						"teamwork":  -1,
						"dreamwork": []bool{false, true},
					}).Debug,
				},
			}

			for _, test := range tests {
				t.Run(test.level, func(t *testing.T) {
					test.f("burger")

					defer b.Reset()

					out := b.String()

					//t.Log(b.String())

					assureSingleNewline(out, t)

					level, file, function, _ := splitMessage(out, t)

					if level != test.level {
						t.Errorf("expected level: '%s'. actual level: '%s'", test.level, level)
					}

					if file != test.file {
						t.Errorf("expected file: '%s'. actual file: '%s'", test.file, file)
					}

					if function != test.function {
						t.Errorf("expected function: '%s'. actual function: '%s'", test.function, function)
					}

					for k, v := range test.fields {
						if ok, fields := hasField(k, v, out, t); !ok {
							t.Errorf("expected fields to contain: '%s=%v. actual fields total: %s", k, v, fields)
						}
					}
				})
			}
		})

		t.Run("Append Fields", func(t *testing.T) {
			defer b.Reset()
			log.InitSimpleLogger(&log.Config{
				Output: b,
			})

			tests := []struct {
				level    string
				file     string
				function string
				fields   log.Fields
				f        func(string)
			}{
				{
					level:    "ERROR",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"one": 1,
					},
					f: log.WithFields(log.Fields{
						"one": 1,
					}).WithFields(log.Fields{}).Error,
				},
				{
					level:    "INFO ",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"sample": "this is a long piece of text",
						"true":   false,
						"false":  true,
					},
					f: log.WithFields(log.Fields{
						"sample": "this is a long piece of text",
					}).WithFields(log.Fields{
						"false": true,
					}).WithFields(log.Fields{
						"true": false,
					}).Info,
				},
				{
					level:    "WARN ",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"one":      nil,
						"okay but": "epic",
					},
					f: log.WithFields(log.Fields{
						"one": nil,
					}).WithFields(log.Fields{
						"okay but": "epic",
					}).Warn,
				},
				{
					level:    "DEBUG",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"teamwork":  -1,
						"dreamwork": []bool{false, true},
					},
					f: log.WithFields(log.Fields{
						"teamwork": -1,
					}).WithFields(log.Fields{
						"dreamwork": []bool{false, true},
					}).Debug,
				},
			}

			for _, test := range tests {
				t.Run(test.level, func(t *testing.T) {
					test.f("burger")

					defer b.Reset()

					out := b.String()

					//t.Log(b.String())

					assureSingleNewline(out, t)

					level, file, function, _ := splitMessage(out, t)

					if level != test.level {
						t.Errorf("expected level: '%s'. actual level: '%s'", test.level, level)
					}

					if file != test.file {
						t.Errorf("expected file: '%s'. actual file: '%s'", test.file, file)
					}

					if function != test.function {
						t.Errorf("expected function: '%s'. actual function: '%s'", test.function, function)
					}

					for k, v := range test.fields {
						if ok, fields := hasField(k, v, out, t); !ok {
							t.Errorf("expected fields to contain: '%s=%v. actual fields total: %s", k, v, fields)
						}
					}
				})
			}
		})

		t.Run("With Error Field", func(t *testing.T) {
			defer b.Reset()
			log.InitSimpleLogger(&log.Config{
				Output: b,
			})

			tests := []struct {
				level    string
				file     string
				function string
				fields   log.Fields
				f        func(string)
			}{
				{
					level:    "ERROR",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"one":   1,
						"error": errors.New("sample text"),
					},
					f: log.WithError(
						errors.New("sample text"),
					).WithFields(log.Fields{
						"one": 1,
					}).Error,
				},
				{
					level:    "INFO ",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"sample": "this is a long piece of text",
						"error":  errors.New("sample text"),
					},
					f: log.WithFields(log.Fields{
						"sample": "this is a long piece of text",
					}).WithError(errors.New("sample text")).Info,
				},
				{
					level:    "WARN ",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"one":   nil,
						"error": errors.New("sample text"),
					},
					f: log.WithFields(log.Fields{
						"one": nil,
					}).WithError(errors.New("sample text")).Warn,
				},
				{
					level:    "DEBUG",
					file:     "log_test.go",
					function: "1()",
					fields: log.Fields{
						"teamwork": -1,
						"error":    errors.New("sample text"),
					},
					f: log.WithFields(log.Fields{
						"teamwork": -1,
					}).WithError(errors.New("sample text")).Debug,
				},
			}

			for _, test := range tests {
				t.Run(test.level, func(t *testing.T) {
					test.f("burger")

					defer b.Reset()

					out := b.String()

					//t.Log(b.String())

					assureSingleNewline(out, t)

					level, file, function, _ := splitMessage(out, t)

					if level != test.level {
						t.Errorf("expected level: '%s'. actual level: '%s'", test.level, level)
					}

					if file != test.file {
						t.Errorf("expected file: '%s'. actual file: '%s'", test.file, file)
					}

					if function != test.function {
						t.Errorf("expected function: '%s'. actual function: '%s'", test.function, function)
					}

					for k, v := range test.fields {
						if ok, fields := hasField(k, v, out, t); !ok {
							t.Errorf("expected fields to contain: '%s=%v. actual fields total: %s", k, v, fields)
						}
					}
				})
			}
		})
	})

	t.Run("LogLevel", func(t *testing.T) {
		tests := []struct {
			levelName string
			level     log.LogLevel
			output    bool
			f         func(string)
		}{
			{
				levelName: "DEBUG",
				level:     log.LogDebug,
				output:    true,
				f:         log.Debug,
			},
			{
				levelName: "ERROR",
				level:     log.LogInformational,
				output:    true,
				f:         log.Error,
			},
			{
				levelName: "INFO ",
				level:     log.LogWarning,
				output:    false,
				f:         log.Info,
			},
			{
				levelName: "WARN ",
				level:     log.LogError,
				output:    false,
				f:         log.Warn,
			},
		}

		var b strings.Builder
		for _, test := range tests {
			t.Run(test.levelName, func(t *testing.T) {
				defer b.Reset()
				log.InitSimpleLogger(&log.Config{
					Output:   &b,
					LogLevel: test.level,
				})

				test.f("sample text")

				if b.Len() > 0 && !test.output {
					t.Errorf("expected no output for log level %d, got '%s'", test.level, b.String())
				}
			})
		}
	})

	t.Run("Clone", func(t *testing.T) {
		defer b.Reset()
		log.InitSimpleLogger(&log.Config{
			Output: b,
		})

		e := log.WithFields(log.Fields{
			"sample": "text",
		})

		e1 := e.Clone().WithFields(log.Fields{
			"fortnite": "borger",
		})

		e = e.WithFields(log.Fields{
			"hello": "world",
		})

		e.Info("e")

		if ok, fields := hasField("fortnite", "borger", b.String(), t); ok {
			t.Errorf("expected to not have '%s=%s' but it did: '%s'", "fortnite", "borger", fields)
		}

		b.Reset()
		e1.Info("e")

		if ok, fields := hasField("hello", "world", b.String(), t); ok {
			t.Errorf("expected to not have '%s=%s' but it did: '%s'", "hello", "world", fields)
		}
	})
}

// Test_JSONLogger is simpler because we only need to test the output, the processing beforehand
// is tested in Test_SimpleLogger
func Test_JSONLogger(t *testing.T) {
	defer b.Reset()

	log.InitJSONLogger(&log.Config{
		Output: b,
	})

	log.WithError(
		errors.New("bepis"),
	).WithFields(log.Fields{
		"hello":  "world",
		"sample": 1,
		"text":   nil,
	}).Error("banana")

	expected := map[string]interface{}{
		"message":   "banana",
		"error":     "bepis",
		"hello":     "world",
		"sample":    float64(1),
		"text":      nil,
		"level":     "ERROR",
		"time":      "<placeholder>",
		"_function": "<placeholder>",
		"_file":     "<placeholder>",
		"_line":     "<placeholder>",
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(b.String()), &data); err != nil {
		t.Fatalf("error unmarshalling buffer: %v", err)
	}

	if len(expected) != len(data) {
		t.Fatalf("expected length: %d. actual length: %d", len(expected), len(data))
	}

	for k, v := range expected {
		val, ok := data[k]
		if !ok {
			t.Errorf("expected '%s' to be in buffer", k)
		}

		// ignore the runtime specific info and timestamp, cant really get that info afaik
		// and checking their presence is good enough
		if !strings.HasPrefix(k, "_") && !(k == "time") {
			if val != v {
				t.Errorf("expected value: %T '%v'. actual value %T '%v'", v, v, val, val)
			}
		}
	}
}
