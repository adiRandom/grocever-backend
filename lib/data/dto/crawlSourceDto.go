package dto

import "fmt"

type CrawlSourceDto struct {
	Url     string
	StoreId int
}

func (dto CrawlSourceDto) String() string {
	return fmt.Sprintf("CrawlSourceDto: (Url: %s StoreId: %d)",
		dto.Url,
		dto.StoreId)
}
