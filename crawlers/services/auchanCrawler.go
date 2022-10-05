package crawlers

import (
	"dealScraper/crawlers/data/constants"
	"dealScraper/crawlers/models"
	"dealScraper/crawlers/utils"
	"github.com/gocolly/colly"
	"strconv"
)

const contentElementQuerySelector = utils.CssSelector(".product-details-info")
const titleElementQuerySelector = utils.CssSelector(".col-md-6 .product-title h1")
const priceElementQuerySelector = utils.CssSelector(".col-md-6 .productDescription .wrapper .price-wrapper-prod-details .big-price #big-price")
const priceAttrib = "data-price"

type AuchanCrawler struct {
}

func (crawler AuchanCrawler) ScrapeProductPage(url string, resCh chan models.CrawlerResult) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML(contentElementQuerySelector.
		String(),
		func(body *colly.HTMLElement) {
			res := models.CrawlerResult{}
			res.ProductName = body.ChildText(titleElementQuerySelector.String())
			price, err := strconv.ParseFloat(body.ChildAttr(priceElementQuerySelector.String(), priceAttrib), 64)

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
