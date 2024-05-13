package app

import (
	"git.snappfood.ir/backend/go/services/bushwack/internal/services"
	"go.uber.org/zap"
)

func (a *application) InitServices(pr *producer, logger *zap.Logger) services.LoggerService {
	return services.NewLoggerService(pr.Redis, a.config, pr.Amqp, logger)
}
