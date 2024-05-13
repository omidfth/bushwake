package app

import (
	"git.snappfood.ir/backend/go/services/bushwack/handler/controllers"
	"git.snappfood.ir/backend/go/services/bushwack/internal/constants/amqpKeys"
)

func (a *application) InitAmqpRouter(pr *producer, c controllers.LoggerController) {
	pr.Amqp.On(amqpKeys.REGISTER, c.Register)
	pr.Amqp.On(amqpKeys.LOG, c.Register)
	pr.Amqp.Serve("bushwake")
}
