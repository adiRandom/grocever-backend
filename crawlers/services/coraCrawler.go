package services

import (
	"crawlers/utils"
	"dealScraper/lib/data/constants"
	"dealScraper/lib/data/models"
	"github.com/gocolly/colly"
	"strconv"
)

const coraContentElementQuerySelector = utils.CssSelector(".product-info-main")
const coraTitleElementQuerySelector = utils.CssSelector(".page-title span[data-ui-id=page-title-wrapper]")
const coraPriceElementQuerySelector = utils.CssSelector(".price-wrapper")
const corePriceAttrib = "data-price-amount"

type CoraCrawler struct {
}

func (crawler CoraCrawler) ScrapeProductPage(url string, resCh chan models.CrawlerResult) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML(coraContentElementQuerySelector.
		String(),
		func(body *colly.HTMLElement) {
			res := models.CrawlerResult{CrawlUrl: url}
			res.ProductName = body.ChildText(coraTitleElementQuerySelector.String())
			price, err := strconv.ParseFloat(body.ChildAttr(coraPriceElementQuerySelector.String(), corePriceAttrib), 64)

			if err != nil {
				return
			}
			res.ProductPrice = price
			res.StoreId = constants.CoraStoreId

			resCh <- res
		})

	err := collyClient.Visit(url)
	if err != nil {
		return
	}
}
