package controllers

import (
	"git.snappfood.ir/backend/go/services/bushwack/internal/repositories/models"
	"git.snappfood.ir/backend/go/services/bushwack/internal/services"
)

type LoggerController interface {
	Register(command models.AmqpModel)
	Log(command models.AmqpModel)
}

type loggerController struct {
	service services.LoggerService
}

func NewLoggerController(service services.LoggerService) LoggerController {
	return &loggerController{service: service}
}

func (c loggerController) Register(command models.AmqpModel) {
	var cmd models.Register
	command.Cast(&cmd)
	c.service.Register(cmd)
}

func (c loggerController) Log(command models.AmqpModel) {
	var cmd models.Log
	command.Cast(&cmd)
	c.service.Log(cmd)
}
