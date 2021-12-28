package logger

import (
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	LUMBERJACK_SCHEME = "lumberjack"
)

func logLevel(in string) (level zapcore.Level, err error) {
	switch in {
	case LEVEL_ERROR:
		level = zap.ErrorLevel
	case LEVEL_WARNING:
		level = zap.WarnLevel
	case LEVEL_INFO:
		level = zap.InfoLevel
	case LEVEL_DEBUG:
		level = zap.DebugLevel
	default:
		err = fmt.Errorf("invalid log level %s", in)
	}
	return level, err
}

func transformLoggerConfig(c *config) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(c.logLevelCore),
		Development: c.development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: c.encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths: c.outputPaths,
	}
}

func registerFileLogger(c *config) {
	if c.logFile == "" {
		return
	}
	zap.RegisterSink(LUMBERJACK_SCHEME, func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename:   c.logFile,
				MaxSize:    c.logFileMaxSize,
				MaxBackups: c.logFileMaxBackups,
				MaxAge:     c.logFileMaxAge,
				Compress:   c.logFileCompress,
			},
		}, nil
	})
	c.outputPaths = append(c.outputPaths, fmt.Sprintf("%s://%s", LUMBERJACK_SCHEME, c.logFile))
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}
