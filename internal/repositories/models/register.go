package models

type Register struct {
	ServiceName string `json:"serviceName"`
	Development bool   `json:"development"`
	OutputPath  string `json:"outputPath"`
	OutputName  string `json:"outputName"`
}

type RegisterCache struct {
	Register
	Token string `json:"token"`
}
