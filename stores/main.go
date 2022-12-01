package main

import (
	"github.com/joho/godotenv"
	"lib/data/database"
	"stores/data/database/entity"
	"stores/gateways/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase(&entity.StoreMetadata{})
	if err != nil {
		return
	}

	c := api.GetClient()
	c.Start()
	println("Started")
}
