![Bushwake Logo](./docs/logo.png "Bushwake")

# Bushwake: Logger Service

Bushwake is a service that facilities logging multi-services.

## How to use
Bushwake uses **`RabbitMQ`** to communication between services.

To use this service, you must Register, first.
```go
amqp.Publish("bushwake",models.AmqpModel{
    Type: "register",
    Body: models.Register{
        ServiceName:     "",
        Development:     false,
        OutputPath:      "",
        OutputName:      "",
	},
})
```

**Token** sent to your service after register completed.

```go
amqp.Publish(command.ServiceName, models.AmqpModel{
			Type: "register",
			Body: token,
		})
```

To log use this code:

```go
amqp.Publish("bushwake", models.AmqpModel{
			Type: "log",
			Body: models.Log{
				Token:     "",
				EventName: "",
				LogType:   "",
				LogFields: nil,
			},
		})
```