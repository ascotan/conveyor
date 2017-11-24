// Package log creates a logging interface using Uber's zap logger
// The following code is based on https://gist.github.com/jinleileiking/95b2efd12563e4abc8d0b8a888b58161
package log

import (
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

var Logger = newZapLogger(false, os.Stdout)

func init() {
	// TODO(jg): reasonable defaults, could override later
	config := Config{
		EncodeLogsAsJson:   false,
		FileLoggingEnabled: true,
		Directory:          "./logs",
		Filename:           "conveyor.log",
		MaxAge:             30,
		MaxSize:            1000,
	}
	Configure(config)
}

// Configuration for logging
type Config struct {
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

// Configure configures the logging system and sets the global logger with an appropriate
// configuration
func Configure(config Config) {
	writers := []zapcore.WriteSyncer{os.Stdout}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	Logger = newZapLogger(config.EncodeLogsAsJson, zapcore.NewMultiWriteSyncer(writers...))
	zap.RedirectStdLog(Logger)
	Logger.Info("logging configured",
		zap.Bool("fileLogging", config.FileLoggingEnabled),
		zap.Bool("jsonLogOutput", config.EncodeLogsAsJson),
		zap.String("logDirectory", config.Directory),
		zap.String("fileName", config.Filename),
		zap.Int("maxSizeMB", config.MaxSize),
		zap.Int("maxBackups", config.MaxBackups),
		zap.Int("maxAgeInDays", config.MaxAge))
}

// newRollingFile takes a logging configuration and uses it to configure
// a rolling log that rolls over every day
func newRollingFile(config Config) zapcore.WriteSyncer {
	if err := os.MkdirAll(config.Directory, 0755); err != nil {
		Logger.Error("failed create log directory", zap.Error(err), zap.String("path", config.Directory))
		return nil
	}

	lj_log := lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,    //megabytes
		MaxAge:     config.MaxAge,     //days
		MaxBackups: config.MaxBackups, //files
		LocalTime:  true,
	}

	c := cron.New()
	// c.AddFunc("* * * * * *", func() { lj_log.Rotate() })
	c.AddFunc("@daily", func() { lj_log.Rotate() })
	c.Start()

	return zapcore.AddSync(&lj_log)
}

// newZapLogger creates a new zap logger. This factory method can either write the output
// as json or for the console (default). It can also take a WriteSyncer interface
// for writing to a specific io.Writer
func newZapLogger(encodeAsJSON bool, output zapcore.WriteSyncer) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "logtime",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	opts := []zap.Option{zap.AddCaller()}
	opts = append(opts, zap.AddStacktrace(zap.WarnLevel))
	encoder := zapcore.NewConsoleEncoder(encCfg)
	if encodeAsJSON {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	return zap.New(zapcore.NewCore(encoder,
		output, zap.NewAtomicLevelAt(zap.DebugLevel)), opts...)
}
