package app

import (
	"git.snappfood.ir/backend/go/services/bushwack/handler/controllers"
	"git.snappfood.ir/backend/go/services/bushwack/internal/services"
)

func (a *application) InitControllers(service services.LoggerService) controllers.LoggerController {
	return controllers.NewLoggerController(service)
}
