package dto

import (
	"dealScraper/crawlers/models"
	"dealScraper/lib/functional"
	"fmt"
)

type ProductProcessDto struct {
	OcrProductDto OcrProductDto          `json:"ocrProduct"`
	CrawlResults  []models.CrawlerResult `json:"crawlResult"`
}

func (dto ProductProcessDto) String() string {
	crawlResultsString := functional.Reduce(dto.CrawlResults,
		func(acc string, crawlResult models.CrawlerResult) string {
			return acc + crawlResult.String() + " "
		}, "")
	return fmt.Sprintf("ProductProcessDto: (OcrProductDto: %s CrawlResults: %s)",
		dto.OcrProductDto.String(),
		crawlResultsString)
}
