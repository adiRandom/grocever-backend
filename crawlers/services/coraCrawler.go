package services

import (
	"crawlers/utils"
	"github.com/gocolly/colly"
	"lib/data/models"
	"lib/data/models/crawl"
	"strconv"
)

const coraContentElementQuerySelector = utils.CssSelector(".product-info-main")
const coraTitleElementQuerySelector = utils.CssSelector(".page-title span[data-ui-id=page-title-wrapper]")
const coraPriceElementQuerySelector = utils.CssSelector(".price-wrapper")
const corePriceAttrib = "data-price-amount"

type CoraCrawler struct {
	store models.StoreMetadata
}

func (crawler CoraCrawler) ScrapeProductPage(url string, resCh chan crawl.ResultModel) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML(coraContentElementQuerySelector.
		String(),
		func(body *colly.HTMLElement) {
			res := crawl.ResultModel{CrawlUrl: url}
			res.ProductName = body.ChildText(coraTitleElementQuerySelector.String())
			price, err := strconv.ParseFloat(body.ChildAttr(coraPriceElementQuerySelector.String(), corePriceAttrib), 32)

			if err != nil {
				return
			}
			res.ProductPrice = float32(price)
			res.Store = crawler.store

			resCh <- res
		})

	err := collyClient.Visit(url)
	if err != nil {
		return
	}
}
