package services

import (
	"crawlers/utils"
	"github.com/gocolly/colly"
	"lib/data/models"
	"lib/data/models/crawl"
	"strconv"
)

const auchanContentElementQuerySelector = utils.CssSelector(".product-details-info")
const auchanTitleElementQuerySelector = utils.CssSelector(".col-md-6 .product-title h1")
const auchanPriceElementQuerySelector = utils.CssSelector(".col-md-6 .productDescription .wrapper .price-wrapper-prod-details .big-price #big-price")
const auchanPriceAttrib = "data-price"

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
			price, err := strconv.ParseFloat(body.ChildAttr(auchanPriceElementQuerySelector.String(), auchanPriceAttrib), 32)

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
