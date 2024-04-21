package app

import (
	"context"
	"flag"
	"git.snappfood.ir/backend/go/services/bushwack/handler/controllers"
	"git.snappfood.ir/backend/go/services/bushwack/internal/constants"
	"git.snappfood.ir/backend/go/services/bushwack/internal/services"
	"git.snappfood.ir/backend/go/services/bushwack/utils"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"log"
)

type Application interface {
	Setup()
	GetContext() context.Context
}

type application struct {
	ctx    context.Context
	config *utils.ServiceConfig
}

func NewApplication(ctx context.Context) Application {
	return &application{ctx: ctx}
}

func (a *application) GetContext() context.Context {
	return a.ctx
}

func (a *application) Setup() {
	path := flag.String("e", constants.DEFAULT_ENV_PATH, "env file path")
	flag.Parse()
	err := a.setupViper(*path)
	if err != nil {
		log.Panic(err.Error())
	}
	app := fx.New(
		fx.Provide(
			a.InitControllers,
			a.InitServices,
			a.InitProducers,
		),
		fx.Invoke(func(pr *producer, logSrv services.LoggerService, c controllers.LoggerController) {
			a.InitAmqpRouter(pr.amqp, c)
			logSrv.Preload()
		}),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
	)
	app.Run()
}
