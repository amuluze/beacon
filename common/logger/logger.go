package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

var once sync.Once

type Logger struct {
	*zap.Logger
	name    string
	lock    sync.Mutex
	loggers map[string]*Logger
}

func init() {
	once.Do(func() {
		defaultLogger = &Logger{
			Logger: zap.New(
				zapcore.NewCore(
					zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
						TimeKey:          "time",
						LevelKey:         "level",
						NameKey:          "logger",
						CallerKey:        "caller",
						MessageKey:       "message",
						StacktraceKey:    "stacktrace",
						LineEnding:       zapcore.DefaultLineEnding,
						EncodeLevel:      cEncodeLevel,
						EncodeTime:       cEncodeTime,
						EncodeDuration:   zapcore.SecondsDurationEncoder,
						EncodeCaller:     cEncodeCaller,
						ConsoleSeparator: " || ",
					}),
					zapcore.AddSync(os.Stdout),
					zap.InfoLevel,
				),
				zap.AddCaller(),
				zap.AddCallerSkip(1),
			),
			name:    "load",
			loggers: make(map[string]*Logger),
		}
	})
}

func (l *Logger) NewLogger(options ...Option) {
	l.lock.Lock()
	defer l.lock.Unlock()
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}

	if _, ok := l.loggers[config.Name]; ok {
		return
	}
	encoder := getEncoder(config)
	writer := getWriter(config)
	level, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}

	newLogger := &Logger{
		Logger: zap.New(
			zapcore.NewCore(encoder, writer, level),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		),
		name:    config.Name,
		loggers: make(map[string]*Logger),
	}
	l.loggers[config.Name] = newLogger
}

func (l *Logger) WithField(fields ...zap.Field) {
	l.Logger = l.Logger.With(fields...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprint(args...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, v...))
}

func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(fmt.Sprint(args...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(fmt.Sprint(args...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(fmt.Sprint(args...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(args ...interface{}) {
	l.Logger.Panic(fmt.Sprint(args...))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panic(fmt.Sprintf(format, v...))
}

// InfoFields 打印带字段的信息日志
func (l *Logger) InfoFields(msg string, fields ...zap.Field) { l.Logger.Info(msg, fields...) }

// DebugFields 打印带字段的调试日志
func (l *Logger) DebugFields(msg string, fields ...zap.Field) { l.Logger.Debug(msg, fields...) }

// WarnFields 打印带字段的警告日志
func (l *Logger) WarnFields(msg string, fields ...zap.Field) { l.Logger.Warn(msg, fields...) }

// ErrorFields 打印带字段的错误日志
func (l *Logger) ErrorFields(msg string, fields ...zap.Field) { l.Logger.Error(msg, fields...) }
