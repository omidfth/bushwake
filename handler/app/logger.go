package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func (a *application) InitLogger() *zap.Logger {
	var encoderCfg zapcore.EncoderConfig
	if a.config.System.DevelopMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       a.config.System.DevelopMode,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
			getPath(a.config.System.LogPath) + "/bushwake" + ".log",
		},
		ErrorOutputPaths: []string{
			"stderr",
			getPath(a.config.System.LogPath) + "/bushwake" + ".log",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

func getPath(logPath string) string {
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return ""
	}
	return logPath
}
