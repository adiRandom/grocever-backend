package main

import (
	"auth/data/entity"
	"auth/gateways/api"
	"lib/microservice"
)

func main() {
	ms := microservice.Microservice{
		HasEnv:     true,
		ApiPortEnv: "API_PORT",
		GetRouter:  api.GetBaseRouter,
		DbEntities: []interface{}{entity.User{}},
	}
	ms.Start()
}
