package services

import (
	"crawlers/utils"
	"fmt"
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
				fmt.Printf("Error crawling %s : %s\n", url, err)
				resCh <- crawl.ResultModel{CrawlUrl: ""}
				return
			}
			res.ProductPrice = float32(price)
			res.Store = crawler.store

			resCh <- res
			fmt.Printf("CoraCrawler from url %s: %v\n", url, res)
		})

	err := collyClient.Visit(url)
	if err != nil {
		fmt.Printf("Error crawling %s : %s\n", url, err)
		resCh <- crawl.ResultModel{CrawlUrl: ""}
		return
	}
}
