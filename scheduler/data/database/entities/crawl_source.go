package entities

import (
	"lib/data/dto/crawl"
	"lib/data/dto/store"
)

type CrawlSource struct {
	Url            string
	StoreId        int    `json:"storeId"`
	StoreName      string `json:"name"`
	OcrHeaderLines int    `json:"ocrHeaderLines"`
	StoreUrl       string `json:"url"`
}

func (s *CrawlSource) ToDto() crawl.SourceDto {
	return crawl.SourceDto{
		Url: s.Url,
		Store: store.MetadataDto{
			StoreId:        s.StoreId,
			Name:           s.StoreName,
			OcrHeaderLines: s.OcrHeaderLines,
			Url:            s.StoreUrl,
		},
	}
}
