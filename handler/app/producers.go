package app

import "git.snappfood.ir/backend/go/services/bushwack/internal/producers"

type producer struct {
	Amqp  producers.AmqpProducer
	Redis producers.RedisProducer
}

func (a *application) InitProducers() *producer {
	redis := producers.NewRedisProducer(a.config.Redis.Host, a.config.Redis.Port)
	amqp := producers.NewAmqpProducer(a.config.Rabbit.ServerName)
	pr := producer{
		Amqp:  amqp,
		Redis: redis,
	}
	return &pr
}
