package microservice

import (
	"github.com/joho/godotenv"
	"lib/api"
	"lib/data/database"
	"os"
)

type Microservice struct {
	HasEnv     bool
	GetRouter  func() *api.Router
	ApiPort    string
	ApiPortEnv string
	DbEntities []interface{}
}

const defaultPort = ":8080"

func (m *Microservice) Start() {
	if m.HasEnv {
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
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

	if m.DbEntities != nil {
		err := database.InitDatabase(m.DbEntities...)
		if err != nil {
			return
		}
	}

	println("Started")
}
