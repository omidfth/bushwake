package app

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (a *application) InitInternalLogger() *zap.Logger {
	var encoderCfg zapcore.EncoderConfig
	if a.config.System.DevelopMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	outputPath := a.config.System.LogPath

	logFile := outputPath + a.config.System.LogName + "-%Y-%m-%d-T%H.log"

	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		panic(err)
	}

	w := zapcore.AddSync(rotator)
	t := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			w,
			zap.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.DebugLevel),
	)

	return zap.New(t)
}
