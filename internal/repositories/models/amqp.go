package models

import "encoding/json"

type AmqpModel struct {
	Type string
	Body interface{}
}

func (a *AmqpModel) Cast(v any) error {
	j, _ := json.Marshal(a.Body)
	return json.Unmarshal(j, &v)
}
