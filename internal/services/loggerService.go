package services

import (
	"encoding/json"
	"git.snappfood.ir/backend/go/services/bushwack/internal/helper"
	"git.snappfood.ir/backend/go/services/bushwack/internal/producers"
	"git.snappfood.ir/backend/go/services/bushwack/internal/repositories/models"
	"git.snappfood.ir/backend/go/services/bushwack/internal/types/logTypes"
	"git.snappfood.ir/backend/go/services/bushwack/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
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

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       command.Development,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
			s.getPath(command.OutputPath) + command.OutputName,
		},
		ErrorOutputPaths: []string{
			"stderr",
			s.getPath(command.ErrorOutputPath) + command.ErrorOutputName,
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	if token == "" {
		token = helper.GetToken()
		j, _ := json.Marshal(models.AmqpModel{
			Type: "register",
			Body: token,
		})
		s.amqp.Publish(command.ServiceName, j)
	}
	return zap.Must(config.Build()), token
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
