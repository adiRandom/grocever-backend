package product

type UserProductDto struct {
	Name                string                  `json:"name"`
	BestPrice           float32                 `json:"bestPrice"`
	PurchaseInstalments []PurchaseInstalmentDto `json:"purchaseInstalments"`
	BestStoreId         uint                    `json:"bestStoreId"`
	BestStoreName       string                  `json:"bestStoreName"`
	BestStoreUrl        string                  `json:"bestStoreUrl"`
}
