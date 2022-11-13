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

const Normal = "normal"
const Prioritized = "prioritized"
const Requeue = "requeue"
