package crawl

import (
	"fmt"
	"lib/data/dto/store"
)

type SourceDto struct {
	Url   string
	Store store.MetadataDto
}

func (dto SourceDto) String() string {
	return fmt.Sprintf("CrawlSourceDto: (Url: %s StoreId: %d)",
		dto.Url,
		dto.Store.StoreId)
}
