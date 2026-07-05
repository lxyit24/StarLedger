package pkg

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"starledger/internal/config"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := config.GlobalConfig.Log

	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var writer zapcore.WriteSyncer
	if cfg.Output == "file" && cfg.FilePath != "" {
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic("failed to open log file: " + err.Error())
		}
		writer = file
	} else {
		writer = os.Stdout
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		level,
	)

	Logger = zap.New(core, zap.AddCaller())
}
