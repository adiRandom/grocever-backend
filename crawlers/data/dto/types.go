package dto

type FreshfulDto struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type MegaImageDto struct {
	Data struct {
		ProductDetails struct {
			Name  string `json:"name"`
			Price struct {
				Value float64 `json:"value"`
			} `json:"price"`
		} `json:"productDetails"`
	} `json:"data"`
}
