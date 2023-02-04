package scheduling

import (
	crawlDto "lib/data/dto/crawl"
	"lib/data/dto/product"
)

type CrawlScheduleDto struct {
	Product crawlDto.ProductDto `json:"product"`
	Type    string              `json:"type"`
}

func NewCrawlScheduleDto(
	product product.PurchaseInstalmentWithUserDto,
	crawlSources []crawlDto.SourceDto,
	crawlType string,
) CrawlScheduleDto {
	return CrawlScheduleDto{
		Product: crawlDto.ProductDto{
			OcrProduct:   product,
			CrawlSources: crawlSources,
		},
		Type: crawlType,
	}
}

func NewRequeueCrawlScheduleDto(
	ocrProductName string,
	crawlSources []crawlDto.SourceDto,
	crawlType string,
) CrawlScheduleDto {
	return CrawlScheduleDto{
		Product: crawlDto.ProductDto{
			OcrProduct: product.PurchaseInstalmentWithUserDto{
				PurchaseInstalmentDto: product.PurchaseInstalmentDto{
					OcrName: ocrProductName,
				},
				UserId: -1,
			},
			CrawlSources: crawlSources,
		},
		Type: crawlType,
	}
}

const Normal = "normal"
const Prioritized = "prioritized"
