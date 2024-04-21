package app

import (
	"git.snappfood.ir/backend/go/services/bushwack/handler/controllers"
	"git.snappfood.ir/backend/go/services/bushwack/internal/constants/amqpKeys"
	"git.snappfood.ir/backend/go/services/bushwack/internal/producers"
)

func (a *application) InitAmqpRouter(pr producers.AmqpProducer, c controllers.LoggerController) {
	pr.On(amqpKeys.REGISTER, c.Register)
	pr.On(amqpKeys.LOG, c.Register)
	pr.Serve("bushwake")
}
