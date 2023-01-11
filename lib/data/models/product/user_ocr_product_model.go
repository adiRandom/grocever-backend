package product

import (
	"lib/data/dto/product_processing"
)

type UserOcrProductModel struct {
	Id         uint
	Qty        float32
	Price      float32
	UserId     uint
	OcrProduct OcrProductModel
	UnitPrice  float32
	StoreId    uint
	UnitType   string
}

func NewUserOcrProductModel(id uint, qty float32, price float32, userId uint, ocrProduct OcrProductModel, unitPrice float32, storeId uint, unitType string) *UserOcrProductModel {
	return &UserOcrProductModel{Id: id, Qty: qty, Price: price, UserId: userId, OcrProduct: ocrProduct, UnitPrice: unitPrice, StoreId: storeId, UnitType: unitType}
}

func (m *UserOcrProductModel) ToDto() product_processing.UserOcrProductDto {
	return product_processing.UserOcrProductDto{
		Id:        m.Id,
		OcrName:   m.OcrProduct.OcrProductName,
		Qty:       m.Qty,
		UnitPrice: m.UnitPrice,
		UnitName:  m.UnitType,
		Price:     m.Price,
		BestPrice: m.OcrProduct.BestPrice,
		StoreId:   m.StoreId,
	}
}

//
//func NewUserOcrProductModelsFromProcessingDto(dto dto.ProductProcessDto) []UserOcrProductModel {

//
//	productModels := functional.Map(dto.CrawlResults, func(crawlResult crawl.CrawlerResult) *Model {
//		return NewProductModel(
//			-1,                            // ID
//			dto.OcrProductDto.ProductName, // Name
//			*crawl.NewCrawlLinkModel(-1, crawlResult.CrawlUrl, crawlResult.Store.StoreId, -1), // CrawlLink
//			crawlResult.Store.StoreId,   // StoreId
//			dto.OcrProductDto.UnitPrice, // Price
//			dto.OcrProductDto.UnitType,  // UnityType
//			[]*OcrProductModel{},        // OcrProducts
//		)
//	})
//
//	ocrProductModel := NewOcrProductModel(
//		dto.OcrProductDto.ProductName,
//		bestPrice,
//		productModels,
//		[]*OcrProductModel{},
//	)
//
//	for _, productModel := range productModels {
//		productModel.OcrProducts = append(productModel.OcrProducts, ocrProductModel)
//	}
//
//	return functional.Map(productModels, func(productModel *Model) UserOcrProductModel {
//		return UserOcrProductModel{
//			Qty:        dto.OcrProductDto.Qty,
//			Price:      dto.OcrProductDto.Price,
//			UserId:     dto.UserId,
//			OcrProduct: *ocrProductModel,
//			Product:    *productModel,
//			Store:      models.NewStoreMetadataFromDto(dto.OcrProductDto.Store),
//		}
//	})
//
//}
