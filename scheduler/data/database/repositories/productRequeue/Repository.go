package productRequeue

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lib/data/database"
	"lib/data/database/repositories"
	"scheduler/data/database/entities"
)

type Repository struct {
	repositories.Repository[entities.ProductRequeueEntity]
}

var pr *Repository = nil

func GetRepository() *Repository {
	if pr == nil {
		pr = &Repository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		pr.Db = db
	}
	return pr
}

func (r *Repository) GetProductsForRequeue() ([]entities.ProductRequeueEntity, error) {
	var products []entities.ProductRequeueEntity
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&products).Error
		if err != nil {
			return err
		}

		return tx.Exec("DELETE FROM product_requeue_entities").Error
	})

	if err != nil {
		return nil, err
	}
	return products, err
}
