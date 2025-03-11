package log

import (
	log "github.com/faelmori/logz"
	logz "github.com/faelmori/logz/logger"
)

var loggerCore logz.LogzCore
var logger logz.LogzLogger

func init() {
	if logger == nil {
		lgz := log.NewLogger("kubero-cli")
		if lgz == nil {
			logger = logz.NewLogger("kubero-cli")
		}
		logger = lgz
	}
	if loggerCore == nil {
		loggerCore = logz.NewLogger("kubero-cli-core")
	}
}

// Logger returns the logger instance
func Logger() log.Logger { return logger }

// SetLogger sets the logger instance
func SetLogger(l log.Logger) { logger = l }

// Fatal logs a message at level Fatal on the standard logger and then calls os.Exit(1).
func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

// Fatalf logs a formatted message at level Fatal on the standard logger and then calls os.Exit(1).
func Fatalf(format string, args ...interface{}) { logger.Fatalf(format, args...) }

// Fatalln logs a message at level Fatal on the standard logger and then calls os.Exit(1).
func Fatalln(args ...interface{}) { logger.Fatalln(args...) }

// Panic logs a message at level Panic on the standard logger and then calls panic with the message.
func Panic(args ...interface{}) { logger.Panic(args...) }

// Panicf logs a formatted message at level Panic on the standard logger and then calls panic with the message.
func Panicf(format string, args ...interface{}) { logger.Panicf(format, args...) }

// Panicln logs a message at level Panic on the standard logger and then calls panic with the message.
func Panicln(args ...interface{}) { logger.Panic(args...) }

// Print calls Output to print to the standard logger.
func Print(args ...interface{}) { logger.Print(args...) }

// Printf calls Output to print to the standard logger.
func Printf(format string, args ...interface{}) { logger.Printf(format, args...) }

// Println calls Output to print to the standard logger.
func Println(args ...interface{}) { logger.Println(args...) }

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger.Debug("Debug", map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Debug", "args": args})
}

// Debugf logs a formatted message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logger.Debug(format, map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Debug", "args": args})
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) { Debug(args...) }

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logger.Info("Info", map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Info", "args": args})
}

// Infof logs a formatted message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger.Info(format, map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Info", "args": args})
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) { Info(args...) }

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger.Warn("Warn", map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Warn", "args": args})
}

// Warnf logs a formatted message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger.Warn(format, map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Warn", "args": args})
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) { Warn(args...) }

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logger.Error("Error", map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Error", "args": args})
}

// Errorf logs a formatted message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger.Error(format, map[string]interface{}{"context": "kubero-cli", "pkg": "log", "method": "Error", "args": args})
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) { Error(args...) }

// SetLevel sets the log level for the logger
func SetLevel(level string) { logger.SetLevel(log.LogLevel(level)) }

// SetFormat sets the log format for the logger
func SetFormat(format string) {
	cfg := logger.GetConfig()
	cfg.SetFormat(log.LogFormat(format))
	logger.SetConfig(cfg)
}

// SetOutput sets the log output for the logger
func SetOutput(output string) {
	cfg := logger.GetConfig()
	cfg.SetOutput(output)
	logger.SetConfig(cfg)
}
