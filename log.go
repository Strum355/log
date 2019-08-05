package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogLevel uint8

// Error Levels that can be used to differentiate logged messages and also
// set the verbosity of logs to display.
const (
	LogDebug LogLevel = iota
	LogInformational
	LogWarning
	LogError
)

type logger interface {
	createLogPoint(log logPoint)
}

type logPoint struct {
	b        *strings.Builder
	level    LogLevel
	fileLine int
	file     string
	funcName string
	msg      string
	fields   Fields
	time     time.Time
}

type Config struct {
	ErrorPrefix string
	WarnPrefix  string
	InfoPrefix  string
	DebugPrefix string
	LogLevel    LogLevel
	Output      io.Writer
	// Will print error level to StdErr
	// UseStdErr is ignored if Output != os.Stdout
	UseStdErr    bool
	logger       logger
	levelPadding int
}

var config *Config

func InitJSONLogger(conf *Config) {
	setDefaults(conf)
	config = conf
	config.logger = newJsonLogger()
}

func InitSimpleLogger(conf *Config) {
	setDefaults(conf)
	setLevelPadding(conf)
	config = conf
	config.logger = newSimpleLogger()

}

func setDefaults(conf *Config) {
	if conf == nil {
		conf = new(Config)
	}

	if conf.LogLevel > LogError {
		panic(fmt.Sprintf("invalid log level %d", conf.LogLevel))
	}

	if conf.ErrorPrefix == "" {
		conf.ErrorPrefix = "ERROR"
	}

	if conf.WarnPrefix == "" {
		conf.WarnPrefix = "WARN"
	}
	if conf.InfoPrefix == "" {
		conf.InfoPrefix = "INFO"
	}
	if conf.DebugPrefix == "" {
		conf.DebugPrefix = "DEBUG"
	}

	if conf.UseStdErr && conf.Output != os.Stdout {
		conf.UseStdErr = false
	}

	if conf.Output == nil {
		conf.Output = os.Stdout
	}
}

func setLevelPadding(conf *Config) {
	maxPadding := func(y int) int {
		if conf.levelPadding > y {
			return conf.levelPadding
		}
		return y
	}

	conf.levelPadding = maxPadding(len(conf.ErrorPrefix))
	conf.levelPadding = maxPadding(len(conf.WarnPrefix))
	conf.levelPadding = maxPadding(len(conf.InfoPrefix))
	conf.levelPadding = maxPadding(len(conf.DebugPrefix))
}
