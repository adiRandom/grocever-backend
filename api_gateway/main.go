package main

import (
	"api_gateway/gateways/api"
	"lib/microservice"
)

func main() {
	ms := microservice.Microservice{
		HasEnv:     true,
		ApiPortEnv: "API_PORT",
		GetRouter:  api.GetBaseRouter,
	}
	ms.Start()
}
