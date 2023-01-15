package scheduling

import (
	"lib/data/dto"
	"lib/data/dto/product"
)

type CrawlScheduleDto struct {
	Product dto.CrawlProductDto `json:"product"`
	Type    string              `json:"type"`
}

func NewCrawlScheduleDto(
	product product.UserOcrProductDto,
	crawlSources []dto.CrawlSourceDto,
	crawlType string,
) CrawlScheduleDto {
	return CrawlScheduleDto{
		Product: dto.CrawlProductDto{
			OcrProduct:   product,
			CrawlSources: crawlSources,
		},
		Type: crawlType,
	}
}

func NewRequeueCrawlScheduleDto(
	ocrProductName string,
	crawlSources []dto.CrawlSourceDto,
	crawlType string,
) CrawlScheduleDto {
	return CrawlScheduleDto{
		Product: dto.CrawlProductDto{
			OcrProduct:   product.UserOcrProductDto{OcrName: ocrProductName, UserId: -1},
			CrawlSources: crawlSources,
		},
		Type: crawlType,
	}
}

const Normal = "normal"
const Prioritized = "prioritized"
