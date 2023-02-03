package services

import (
	"crawlers/utils"
	"fmt"
	"github.com/gocolly/colly"
	"lib/data/models"
	"lib/data/models/crawl"
	"strconv"
	"strings"
)

const megaContentElementQuerySelector = utils.CssSelector("[data-testid=product-details-section]")
const megaTitleElementQuerySelector = utils.CssSelector("[data-testid=product-common-header-title]")
const megaPriceElementQuerySelector = utils.CssSelector("[data-testid=product-block-price]")

type MegaImageCrawler struct {
	store models.StoreMetadata
}

func (crawler MegaImageCrawler) ScrapeProductPage(url string, resCh chan crawl.ResultModel) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML("body", func(body *colly.HTMLElement) {
		print(body.DOM.Html())
		if body.DOM.Find(megaContentElementQuerySelector.String()).Length() == 0 {
			// Not url to product page
			resCh <- crawl.ResultModel{CrawlUrl: ""}
		}
	})

	collyClient.OnHTML(megaContentElementQuerySelector.
		String(),
		func(body *colly.HTMLElement) {
			res := crawl.ResultModel{CrawlUrl: url}
			res.ProductName = body.ChildText(megaTitleElementQuerySelector.String())
			priceWithCurrency := body.ChildText(megaPriceElementQuerySelector.String())
			priceString := strings.Replace(priceWithCurrency[0:len(priceWithCurrency)-4], ",", ".", 1)
			price, err := strconv.ParseFloat(priceString, 32)

			if err != nil {
				fmt.Printf("Error crawling %s : %s\n", url, err)
				resCh <- crawl.ResultModel{CrawlUrl: ""}
				return
			}
			res.ProductPrice = float32(price)
			res.Store = crawler.store

			resCh <- res
			fmt.Printf("Mega Image from url %s: %s - %s - %f\n", url, res.ProductName, res.ProductPrice, res.ProductPrice)
		})

	err := collyClient.Visit(url)
	if err != nil {
		fmt.Printf("Error crawling %s : %s\n", url, err)
		resCh <- crawl.ResultModel{CrawlUrl: ""}
		return
	}
}
