package dto

type FreshfulDto struct {
	Name   string  `json:"name"`
	Price  float32 `json:"price"`
	Images []struct {
		Thumbnail struct {
			Default string `json:"default"`
		} `json:"thumbnail"`
		Large struct {
			Default string `json:"default"`
		} `json:"large"`
		Extralarge struct {
			Default string `json:"default"`
		} `json:"extralarge"`
	} `json:"images"`
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
