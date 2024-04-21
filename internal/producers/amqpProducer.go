package producers

import (
	"context"
	"encoding/json"
	"fmt"
	"git.snappfood.ir/backend/go/services/bushwack/internal/repositories/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

type AmqpProducer interface {
	Publish(key string, body []byte)
	Serve(key string)
	On(key string, f func(data models.AmqpModel)) *amqpRoute
}

type amqpRoute struct {
	handler amqpHandler
}

func (s *amqpProducer) On(key string, f func(data models.AmqpModel)) *amqpRoute {
	return s.addAmqpHandler(key, amqpFunc(f))
}

type amqpHandler interface {
	Serve(models.AmqpModel)
}

type amqpFunc func(models.AmqpModel)

func (f amqpFunc) Serve(m models.AmqpModel) {
	f(m)
}

type amqpProducer struct {
	conn   *amqp.Connection
	events map[string]*amqpRoute
}

func (s *amqpProducer) addAmqpHandler(key string, handler amqpHandler) *amqpRoute {
	route := amqpRoute{handler: handler}
	s.events[key] = &route
	return &route
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err.Error())
	}
}

func NewAmqpProducer(amqpServer string) AmqpProducer {
	var counts int64
	var backOff = 1 * time.Second
	var conn *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			conn = c
			break
		}

		if counts > 5 {
			failOnError(err, "Failed to connect to RabbitMQ")
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"task",  // name
		"topic", // type
		false,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	return &amqpProducer{conn: conn, events: make(map[string]*amqpRoute)}
}

func (s *amqpProducer) Publish(key string, body []byte) {
	ch, err := s.conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", string(body))
}

func (s *amqpProducer) Serve(key string) {
	go func() {
		ch, err := s.conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			key,   // name
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		failOnError(err, "Failed to declare a queue")

		err = ch.QueueBind(q.Name, "topic", "task", false, nil)
		failOnError(err, "Failed to bind a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")
		var forever chan struct{}

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				var data models.AmqpModel
				err := json.Unmarshal(d.Body, &data)
				if err != nil {
					continue
				}
				val, ok := s.events[data.Type]
				if !ok {
					continue
				}
				val.handler.Serve(data)
			}
		}()
		<-forever
	}()
}
