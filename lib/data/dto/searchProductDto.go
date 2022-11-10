package dto

import "lib/functional"

type SearchProductDto struct {
	OcrProduct   OcrProductDto    `json:"ocrProduct"`
	CrawlSources []CrawlSourceDto `json:"crawlSources"`
}

func (dto SearchProductDto) String() string {
	crawlSourcesString := functional.Reduce(dto.CrawlSources,
		func(acc string, crawlSource CrawlSourceDto) string {
			return acc + crawlSource.String() + " "
		}, "")
	return "SearchProductDto: (OcrProduct: " +
		dto.OcrProduct.String() +
		" CrawlSources: " + crawlSourcesString +
		")"
}
