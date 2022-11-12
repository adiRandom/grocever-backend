package dto

type FreshfulDto struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type MegaImageDto struct {
	Data struct {
		ProductDetails struct {
			Name  string `json:"name"`
			Price struct {
				Value float32 `json:"value"`
			} `json:"price"`
		} `json:"productDetails"`
	} `json:"data"`
}
