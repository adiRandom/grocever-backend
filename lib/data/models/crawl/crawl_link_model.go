package crawl

type LinkModel struct {
	Id        int
	Url       string
	StoreId   int
	ProductId int
}

func NewCrawlLinkModel(id int, url string, storeId int, productId int) *LinkModel {
	return &LinkModel{Id: id, Url: url, StoreId: storeId, ProductId: productId}
}
