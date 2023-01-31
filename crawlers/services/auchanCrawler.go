package services

import (
	"crawlers/utils"
	"fmt"
	"github.com/gocolly/colly"
	"lib/data/models"
	"lib/data/models/crawl"
	"strconv"
)

const auchanContentElementQuerySelector = utils.CssSelector(".vtex-flex-layout-0-x-flexCol--prodDetailsDesktop")
const auchanTitleElementQuerySelector = utils.CssSelector(".vtex-store-components-3-x-productBrand")
const auchanIntPriceElementQuerySelector = utils.CssSelector(".vtex-product-price-1-x-currencyInteger")
const auchanDecimalPriceElementQuerySelector = utils.CssSelector(".vtex-product-price-1-x-currencyFraction")

type AuchanCrawler struct {
	store models.StoreMetadata
}

func (crawler AuchanCrawler) ScrapeProductPage(url string, resCh chan crawl.ResultModel) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML(auchanContentElementQuerySelector.
		String(),
		func(body *colly.HTMLElement) {
			res := crawl.ResultModel{CrawlUrl: url}
			res.ProductName = body.ChildText(auchanTitleElementQuerySelector.String())
			priceString := body.ChildText(auchanIntPriceElementQuerySelector.String()) + "." + body.ChildText(auchanDecimalPriceElementQuerySelector.String())
			price, err := strconv.ParseFloat(priceString, 32)

			if err != nil {
				fmt.Printf("Error crawling %s : %s\n", url, err)
				resCh <- crawl.ResultModel{CrawlUrl: ""}
				return
			}
			res.ProductPrice = float32(price)
			res.Store = crawler.store

			resCh <- res
			fmt.Printf("AuchanCrawler from url %s: %v\n", url, res)
		})

	err := collyClient.Visit(url)
	if err != nil {
		fmt.Printf("Error crawling %s : %s\n", url, err)
		resCh <- crawl.ResultModel{CrawlUrl: ""}
		return
	}
}
