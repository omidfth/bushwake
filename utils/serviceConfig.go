package utils

type ServiceConfig struct {
	System system `mapstructure:"system"`
	Rabbit rabbit `mapstructure:"rabbit"`
	Redis  redis  `mapstructure:"redis"`
}

type system struct {
	LogPath     string `mapstructure:"log_path"`
	DevelopMode bool   `mapstructure:"develop_mode"`
}
type rabbit struct {
	ServerName   string `mapstructure:"server_name"`
	ExchangeName string `mapstructure:"exchange_name"`
	PublishKind  string `mapstructure:"publish_kind"`
}

type redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	TTL  int    `mapstructure:"ttl"`
}
