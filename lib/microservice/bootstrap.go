package microservice

import (
	"context"
	"github.com/joho/godotenv"
	"lib/data/database"
	"os"
)

func (m *Microservice) Start() {
	if m.HasEnv {
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	if m.DbEntities != nil {
		err := database.InitDatabase(m.DbEntities...)
		if err != nil {
			return
		}
	}

	router := m.GetRouter()

	if router != nil {
		port := m.ApiPort
		apiPortEnv := os.Getenv(m.ApiPortEnv)
		if port == "" {
			if apiPortEnv != "" {
				port = apiPortEnv
			} else {
				port = defaultPort
			}
		}
		router.Run(port)
	}

	println("Started")
}

func (m *AsyncMicroservice[T]) Start() {
	m.Microservice.Start()

	go m.CreateMessageBroker().Start(context.Background())
}
