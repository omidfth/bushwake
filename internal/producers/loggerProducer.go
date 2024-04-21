package producers

import (
	"git.snappfood.ir/backend/go/services/bushwack/internal/repositories/models"
	"git.snappfood.ir/backend/go/services/bushwack/internal/types/logTypes"
)

type LoggerProducer interface {
	Serve(log models.Log)
	On(key logTypes.LogType, f loggerFunc) *loggerRoute
}

type loggerRoute struct {
	handler loggerHandler
}

func (s *loggerProducer) On(key logTypes.LogType, f loggerFunc) *loggerRoute {
	return s.addLoggerHandler(key, f)
}

type loggerHandler interface {
	Serve(logModel models.Log)
}

type loggerFunc func(logModel models.Log)

func (f loggerFunc) Serve(logModel models.Log) {
	f(logModel)
}

type loggerProducer struct {
	events map[logTypes.LogType]*loggerRoute
}

func (s *loggerProducer) addLoggerHandler(key logTypes.LogType, handler loggerHandler) *loggerRoute {
	route := loggerRoute{handler: handler}
	s.events[key] = &route
	return &route
}

func NewLoggerProducer() LoggerProducer {
	return &loggerProducer{events: make(map[logTypes.LogType]*loggerRoute)}
}

func (s *loggerProducer) Serve(log models.Log) {
	val, ok := s.events[log.LogType]
	if !ok {
		return
	}
	val.handler.Serve(log)
}
