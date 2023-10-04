// Package azlog provides a thin wrapper around the zap logging package to
// enable structured logging on Azure. This is still a work in progress and
// is untested.
package azlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Local       bool
	MinLogLevel zapcore.Level
}

func New(cfg *Config) (*zap.Logger, error) {
	if cfg.Local {
		zCfg := zap.NewDevelopmentConfig()
		zCfg.Level = zap.NewAtomicLevelAt(cfg.MinLogLevel)
		return zCfg.Build()
	}

	zCfg := &zap.Config{
		Level:       zap.NewAtomicLevelAt(cfg.MinLogLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return zCfg.Build(zap.AddStacktrace(zap.DPanicLevel))
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	/*
		TODO: Figure out what the appropriate logging method is here.
		See https://learn.microsoft.com/en-us/dotnet/api/microsoft.extensions.logging.loglevel?view=dotnet-plat-ext-7.0
			Critical 	5 	 Logs that describe an unrecoverable application or system crash, or a catastrophic failure that requires immediate attention.
			Debug 	1 	 Logs that are used for interactive investigation during development. These logs should primarily contain information useful for debugging and have no long-term value.
			Error 	4 	 Logs that highlight when the current flow of execution is stopped due to a failure. These should indicate a failure in the current activity, not an application-wide failure.
			Information 	2 	 Logs that track the general flow of the application. These logs should have long-term value.
			None 	6 	 Not used for writing log messages. Specifies that a logging category should not write any messages.
			Trace 	0 	 Logs that contain the most detailed messages. These messages may contain sensitive application data. These messages are disabled by default and should never be enabled in a production environment.
			Warning 	3 	 Logs that highlight an abnormal or unexpected event in the application flow, but do not otherwise cause the application execution to stop.
	*/
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARNING")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("CRITICAL")
	case zapcore.PanicLevel:
		enc.AppendString("ALERT")
	case zapcore.FatalLevel:
		enc.AppendString("EMERGENCY")
	}
}
