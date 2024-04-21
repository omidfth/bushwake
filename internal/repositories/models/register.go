package models

type Register struct {
	ServiceName     string `json:"serviceName"`
	Development     bool   `json:"development"`
	OutputPath      string `json:"outputPath"`
	OutputName      string `json:"outputName"`
	ErrorOutputPath string `json:"errorOutputPath"`
	ErrorOutputName string `json:"errorOutputName"`
}

type RegisterCache struct {
	Register
	Token string `json:"token"`
}
