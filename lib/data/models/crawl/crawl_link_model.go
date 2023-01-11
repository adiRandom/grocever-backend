package crawl

import "lib/data/dto"

type LinkModel struct {
	Id        uint
	Url       string
	StoreId   int
	ProductId uint
}

func NewCrawlLinkModel(id uint, url string, storeId int, productId uint) *LinkModel {
	return &LinkModel{Id: id, Url: url, StoreId: storeId, ProductId: productId}
}

func (model *LinkModel) ToCrawlSourceDto() dto.CrawlSourceDto {
	return dto.CrawlSourceDto{
		Url:     model.Url,
		StoreId: model.StoreId,
	}
}
