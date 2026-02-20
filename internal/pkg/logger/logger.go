package logger

import (
	"fmt"
	"os"
	"time"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logLevel - miniman levl of logs
func NewLogger(logLevel string) (*zap.Logger, func() error, error) {

	// create managable level of logging
	// managable here means that we can change level of logging in any moment
	lvl := zap.NewAtomicLevel()

	// convert logLevel string (can be warn, debug, info and etc) into zapcore.Level
	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, nil, fmt.Errorf("Unmarshal log level: %w", err)
	}

	// first argument of mkdir func - path
	if err := os.Mkdir("logs", 0755); err != nil {
		return nil, nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	////////////////////////////////////////////////////////
	// make graceful entry of data
	timestamp := time.Now().UTC().Format("2026-01-02T15-04-05.000000")
	// add this data into file with logs (which was created before)
	logFilePath := filepath.Join("logs", fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file %w", err)
	}

	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, logFile.Close, nil
}