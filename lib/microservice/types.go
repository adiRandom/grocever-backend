package microservice

import (
	"lib/api"
	"lib/events/rabbitmq"
)

type Microservice struct {
	HasEnv     bool
	GetRouter  func() *api.Router
	ApiPort    string
	ApiPortEnv string
	DbEntities []interface{}
}

const defaultPort = ":8080"

type AsyncMicroservice[T any] struct {
	Microservice
	MessageBroker rabbitmq.JsonBroker[T]
}
