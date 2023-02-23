package services

import (
	"fmt"
	"lib/data/constants"
	crawl2 "lib/data/dto/crawl"
	"lib/data/models"
	"lib/data/models/crawl"
	"lib/functional"
)

func getCrawler(store models.StoreMetadata) Crawler {
	switch store.StoreId {
	case constants.AuchanStoreId:
		{
			return &AuchanCrawler{store}
		}
	case constants.FreshfulStoreId:
		{
			return &FreshfulCrawler{store}
		}
	case constants.CoraStoreId:
		{
			return &CoraCrawler{store}
		}
	}
	return nil
}

func crawlProductPage(src crawl2.SourceDto, resCh chan crawl.ResultModel) {
	store := models.NewStoreMetadataFromDto(src.Store)
	crawler := getCrawler(store)
	if crawler == nil {
		return
	}
	crawler.ScrapeProductPage(src.Url, resCh)
}

func crawlProductPages(srcs []crawl2.SourceDto, resCh chan crawl.ResultModel) {
	for _, src := range srcs {
		go crawlProductPage(src, resCh)
	}
}

func CrawlProductPages(srcs []crawl2.SourceDto) []crawl.ResultModel {
	resCh := make(chan crawl.ResultModel)

	filteredSrcs := functional.Filter(srcs, func(src crawl2.SourceDto) bool {
		return src.Url != ""
	})

	crawlProductPages(filteredSrcs, resCh)

	var res []crawl.ResultModel
	for range filteredSrcs {
		res = append(res, <-resCh)
		fmt.Printf("Got result: %+v\n", res)
	}

	return functional.Filter(res, func(r crawl.ResultModel) bool {
		return r.CrawlUrl != "" && r.ProductPrice != 0 && r.ProductName != ""
	})
}
