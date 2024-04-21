package models

import (
	"git.snappfood.ir/backend/go/services/bushwack/internal/types/logTypes"
	"go.uber.org/zap"
)

type Log struct {
	Token     string           `json:"token"`
	EventName string           `json:"event_name"`
	LogType   logTypes.LogType `json:"log_type"`
	LogFields []LogFields      `json:"log_fields"`
}

func (m *Log) GetFields() []zap.Field {
	var fields []zap.Field
	for _, f := range m.LogFields {
		switch val := f.Value.(type) {
		case string:
			fields = append(fields, zap.String(f.Key, val))
		case int:
			fields = append(fields, zap.Int(f.Key, val))
		case interface{}:
			fields = append(fields, zap.Any(f.Key, val))
		}
	}

	return fields
}

type LogFields struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
