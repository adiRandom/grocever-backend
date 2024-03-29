package product

import (
	"lib/data/dto"
	"lib/data/models/crawl"
)

type Model struct {
	ID          int
	Name        string
	ImageUrl    string
	CrawlLink   crawl.LinkModel
	StoreId     int
	Price       float32
	UnityType   string
	OcrProducts []*OcrProductModel
}

func NewProductModel(
	ID int,
	name string,
	imageUrl string,
	crawlLink crawl.LinkModel,
	storeId int,
	price float32,
	unityType string,
	ocrProducts []*OcrProductModel) *Model {
	return &Model{ID: ID, Name: name, ImageUrl: imageUrl, CrawlLink: crawlLink, StoreId: storeId, Price: price, UnityType: unityType, OcrProducts: ocrProducts}
}

func NewProductModelsFromProcessDto(dto dto.ProductProcessDto, existingOcrProduct *OcrProductModel) []*Model {
	productModels := make([]*Model, len(dto.CrawlResults))
	for i, crawlResult := range dto.CrawlResults {
		productModels[i] = NewProductModel(
			-1,
			crawlResult.ProductName,
			crawlResult.ImageUrl,
			crawl.LinkModel{
				Id:        -1,
				Url:       crawlResult.CrawlUrl,
				StoreId:   crawlResult.Store.StoreId,
				ProductId: -1,
			},
			crawlResult.Store.StoreId,
			crawlResult.ProductPrice,
			dto.OcrProduct.UnitName,
			[]*OcrProductModel{},
		)
	}

	if existingOcrProduct == nil {
		ocrProduct := NewOcrProductModel(dto.OcrProduct.OcrName, nil, productModels, []*OcrProductModel{})

		for _, productModel := range productModels {
			productModel.OcrProducts = append(productModel.OcrProducts, ocrProduct)
		}

	} else {
		for _, productModel := range productModels {
			productModel.OcrProducts = append(productModel.OcrProducts, existingOcrProduct)
		}
	}

	return productModels

}
