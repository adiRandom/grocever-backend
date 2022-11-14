package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	"scheduler/data/database/entities"
)

type ProductRequeueRepository struct {
	repositories.Repository[entities.ProductRequeueEntity]
}

var pr *ProductRequeueRepository = nil

func GetProductRequeueRepository() *ProductRequeueRepository {
	if pr == nil {
		pr = &ProductRequeueRepository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		pr.Db = db
	}
	return pr
}
