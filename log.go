package log

import (
	"io"
	"os"

	"github.com/bwmarrin/lit"
)

type Config struct {
	ErrorPrefix string
	WarnPrefix  string
	InfoPrefix  string
	DebugPrefix string
	Output      io.Writer
	// Will print error level to StdErr
	// UseStdErr is ignored if Output != os.Stdout
	UseStdErr bool
}

var config Config

func Init(conf Config) {
	if conf.ErrorPrefix == "" {
		conf.ErrorPrefix = "ERROR"
	}
	if conf.WarnPrefix == "" {
		conf.WarnPrefix = "WARN "
	}
	if conf.InfoPrefix == "" {
		conf.InfoPrefix = "INFO "
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

	config = conf

	lit.Prefix = ""
	lit.PrefixError = conf.ErrorPrefix
	lit.PrefixDebug = conf.DebugPrefix
	lit.PrefixWarning = conf.WarnPrefix
	lit.PrefixInformational = conf.InfoPrefix
	lit.LogLevel = lit.LogDebug
}
