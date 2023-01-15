package store

type MetadataDto struct {
	StoreId        int    `json:"storeId"`
	Name           string `json:"name"`
	OcrHeaderLines int    `json:"ocrHeaderLines"`
	Url            string `json:"url"`
}
