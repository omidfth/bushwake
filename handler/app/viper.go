package app

import (
	"errors"
	"fmt"
	"git.snappfood.ir/backend/go/services/bushwack/utils"
	"github.com/spf13/viper"
	"os"
)

func (a *application) setupViper(path string) error {
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	f, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("cannot read config file: %s", err.Error())
		return errors.New(msg)
	}
	err = viper.ReadConfig(f)
	if err != nil {
		msg := fmt.Sprintf("viper read config error: %s", err.Error())
		return errors.New(msg)
	}
	var c utils.ServiceConfig
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}
	a.config = &c
	return nil
}
