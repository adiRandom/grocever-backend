package repositories

import "productProcessing/data/database/entities"

type bestPriceResult struct {
	ocrName string
	product *entities.ProductEntity
}

type ocrNameAndProductId struct {
	ocrName   string
	productId uint
}

type ocrProductsLinksDenied map[ocrNameAndProductId]struct{}

func (m ocrProductsLinksDenied) IsLinkDenied(ocrName string, productId uint) bool {
	_, ok := m[ocrNameAndProductId{ocrName: ocrName, productId: productId}]
	return ok
}
