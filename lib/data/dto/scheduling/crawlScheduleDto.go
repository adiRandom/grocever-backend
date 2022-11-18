package scheduling

import (
	"fmt"
	"lib/data/dto"
)

type CrawlDto struct {
	Product dto.CrawlProductDto `json:"product"`
	Type    string              `json:"type"`
}

func (dto CrawlDto) String() string {
	return fmt.Sprintf(
		"CrawlDto: (Type: %s, Product: %s)",
		dto.Type,
		dto.Product.String(),
	)
}

func NewCrawlDto(
	product dto.OcrProductDto,
	crawlSource dto.CrawlSourceDto,
	crawlType string,
) CrawlDto {
	return CrawlDto{
		Product: dto.CrawlProductDto{
			OcrProduct:   product,
			CrawlSources: []dto.CrawlSourceDto{crawlSource},
		},
		Type: crawlType,
	}
}

const Normal = "normal"
const Prioritized = "prioritized"
const Requeue = "requeue"
