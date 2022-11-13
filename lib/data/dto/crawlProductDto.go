package dto

import "lib/functional"

type CrawlProductDto struct {
	OcrProduct   OcrProductDto    `json:"ocrProduct"`
	CrawlSources []CrawlSourceDto `json:"crawlSources"`
}

func (dto CrawlProductDto) String() string {
	crawlSourcesString := functional.Reduce(dto.CrawlSources,
		func(acc string, crawlSource CrawlSourceDto) string {
			return acc + crawlSource.String() + " "
		}, "")
	return "CrawlProductDto: (OcrProducts: " +
		dto.OcrProduct.String() +
		" CrawlSources: " + crawlSourcesString +
		")"
}
