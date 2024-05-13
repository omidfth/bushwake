package services

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"git.snappfood.ir/backend/go/services/bushwack/internal/producers"
	"git.snappfood.ir/backend/go/services/bushwack/internal/repositories/models"
	"git.snappfood.ir/backend/go/services/bushwack/internal/types/logTypes"
	"git.snappfood.ir/backend/go/services/bushwack/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerService interface {
	Register(command models.Register)
	Preload()
	Log(log models.Log)
}

type loggerService struct {
	loggers        map[string]producers.LoggerProducer
	redis          producers.RedisProducer
	amqp           producers.AmqpProducer
	loggerPr       producers.LoggerProducer
	config         *utils.ServiceConfig
	internalLogger *zap.Logger
	mu             sync.Mutex
}

func NewLoggerService(
	redis producers.RedisProducer,
	config *utils.ServiceConfig,
	amqp producers.AmqpProducer,
	internalLogger *zap.Logger,
) LoggerService {
	return &loggerService{
		loggers:        make(map[string]producers.LoggerProducer),
		redis:          redis,
		config:         config,
		amqp:           amqp,
		internalLogger: internalLogger,
	}
}

func (s *loggerService) Preload() {
	commands := s.getCache()
	if commands == nil {
		return
	}
	for _, c := range commands {
		logger, token := s.addService(c.Register, c.Token)
		pr := s.addLoggerListeners(logger)
		s.loggers[token] = pr
	}
}

func (s *loggerService) Register(command models.Register) {
	logger, token := s.addService(command, "")
	pr := s.addLoggerListeners(logger)
	s.loggers[token] = pr
}

func (s *loggerService) addLoggerListeners(logger *zap.Logger) producers.LoggerProducer {
	pr := producers.NewLoggerProducer()
	pr.On(logTypes.DEBUG, func(command models.Log) { logger.Debug(command.EventName, command.GetFields()...) })
	pr.On(logTypes.INFO, func(command models.Log) { logger.Info(command.EventName, command.GetFields()...) })
	pr.On(logTypes.WARN, func(command models.Log) { logger.Warn(command.EventName, command.GetFields()...) })
	pr.On(logTypes.ERROR, func(command models.Log) { logger.Error(command.EventName, command.GetFields()...) })
	pr.On(logTypes.PANIC, func(command models.Log) { logger.Panic(command.EventName, command.GetFields()...) })
	pr.On(logTypes.FATAL, func(command models.Log) { logger.Fatal(command.EventName, command.GetFields()...) })
	return pr
}
func (s *loggerService) Log(log models.Log) {
	logger, ok := s.loggers[log.Token]
	if !ok {
		s.internalLogger.Error("token not found", zap.String("token", log.Token))
		return
	}
	logger.Serve(log)
}

func (s *loggerService) addService(command models.Register, token string) (*zap.Logger, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var encoderCfg zapcore.EncoderConfig
	if command.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	outputPath := s.getPath(command.OutputPath)

	logFile := outputPath + command.OutputName + "-%Y-%m-%d-T%H.log"

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

	return zap.New(t), token
}

func (s *loggerService) setCache(command models.Register, token string) {
	registeredServices := s.getCache()
	isAllowSetCache := true
	if registeredServices != nil {
		for _, rs := range registeredServices {
			if rs.ServiceName == command.ServiceName {
				isAllowSetCache = false
				break
			}
		}
	}

	if !isAllowSetCache {
		return
	}
	cache := models.RegisterCache{
		Register: command,
		Token:    token,
	}
	j, _ := json.Marshal(cache)
	s.redis.Set("loggers", string(j), time.Duration(s.config.Redis.TTL)*time.Second)
}

func (s *loggerService) getCache() []models.RegisterCache {
	data := s.redis.Get("loggers")
	if data == "" {
		return nil
	}
	var registeredServices []models.RegisterCache
	json.Unmarshal([]byte(data), &registeredServices)
	return registeredServices
}

func (s *loggerService) getPath(logPath string) string {
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return ""
	}
	return logPath
}
