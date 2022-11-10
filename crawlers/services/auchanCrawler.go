package services

import (
	"crawlers/utils"
	"github.com/gocolly/colly"
	"lib/data/constants"
	"lib/data/models"
	"strconv"
)

const auchanContentElementQuerySelector = utils.CssSelector(".product-details-info")
const auchanTitleElementQuerySelector = utils.CssSelector(".col-md-6 .product-title h1")
const auchanPriceElementQuerySelector = utils.CssSelector(".col-md-6 .productDescription .wrapper .price-wrapper-prod-details .big-price #big-price")
const auchanPriceAttrib = "data-price"

type AuchanCrawler struct {
}

func (crawler AuchanCrawler) ScrapeProductPage(url string, resCh chan models.CrawlerResult) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML(auchanContentElementQuerySelector.
		String(),
		func(body *colly.HTMLElement) {
			res := models.CrawlerResult{CrawlUrl: url}
			res.ProductName = body.ChildText(auchanTitleElementQuerySelector.String())
			price, err := strconv.ParseFloat(body.ChildAttr(auchanPriceElementQuerySelector.String(), auchanPriceAttrib), 64)

			if err != nil {
				return
			}
			res.ProductPrice = price
			res.StoreId = constants.AuchanStoreId

			resCh <- res
		})

	err := collyClient.Visit(url)
	if err != nil {
		return
	}
}
