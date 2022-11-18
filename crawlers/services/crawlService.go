package services

import (
	"fmt"
	"lib/data/constants"
	"lib/data/dto"
	"lib/data/models"
)

func getCrawler(storeId int) Crawler {
	switch storeId {
	case constants.AuchanStoreId:
		{
			return &AuchanCrawler{}
		}
	case constants.FreshfulStoreId:
		{
			return &FreshfulCrawler{}
		}
	case constants.MegaImageStoreId:
		{
			return &MegaImageCrawler{}
		}
	case constants.CoraStoreId:
		{
			return &CoraCrawler{}
		}
	}
	return nil
}

func crawlProductPage(src dto.CrawlSourceDto, resCh chan models.CrawlerResult) {
	crawler := getCrawler(src.StoreId)
	if crawler == nil {
		return
	}
	crawler.ScrapeProductPage(src.Url, resCh)
}

func crawlProductPages(srcs []dto.CrawlSourceDto, resCh chan models.CrawlerResult) {
	for _, src := range srcs {
		go crawlProductPage(src, resCh)
	}
}

func CrawlProductPages(srcs []dto.CrawlSourceDto) []models.CrawlerResult {
	resCh := make(chan models.CrawlerResult)
	crawlProductPages(srcs, resCh)

	var res []models.CrawlerResult
	for range srcs {
		res = append(res, <-resCh)
		fmt.Printf("Got result: %+v\n", res)
	}
	return res
}
