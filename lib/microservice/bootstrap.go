package microservice

import (
	"github.com/joho/godotenv"
	"lib/api"
	"os"
)

type Microservice struct {
	HasEnv     bool
	Router     *api.Router
	ApiPort    string
	ApiPortEnv string
}

const defaultPort = ":8080"

func (m *Microservice) Start() {
	if m.HasEnv {
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	if m.Router != nil {
		port := m.ApiPort
		apiPortEnv := os.Getenv(m.ApiPortEnv)
		if port == "" {
			if apiPortEnv != "" {
				port = apiPortEnv
			} else {
				port = defaultPort
			}
		}
		m.Router.Run(port)
	}

	println("Started")
}
