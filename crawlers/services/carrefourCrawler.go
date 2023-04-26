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

const carrefourContentElementQuerySelector = utils.CssSelector(".product-page-view,.product-page-view-new,.column,.main")
const carrefourTitleElementQuerySelector = utils.CssSelector("[data-ui-id=page-title-wrapper]")
const carrefourIntPriceElementQuerySelector = utils.CssSelector(".price span:nth-child(1)")
const carrefourDecimalPriceElementQuerySelector = utils.CssSelector(".price span:nth-child(2) span:nth-child(2)")
const carrefourImageElementQuerySelector = utils.CssSelector(".gallery-large")

type CarrefourCrawler struct {
	store models.StoreMetadata
}

func (crawler CarrefourCrawler) ScrapeProductPage(url string, resCh chan crawl.ResultModel) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML("body", func(body *colly.HTMLElement) {
		if body.DOM.Find(carrefourContentElementQuerySelector.String()).Length() == 0 {
			// Not url to product page
			resCh <- crawl.ResultModel{CrawlUrl: ""}
			return
		}

		res := crawl.ResultModel{CrawlUrl: url}
		productName, err := body.DOM.Find(carrefourTitleElementQuerySelector.String()).Html()

		if err == nil {
			res.ProductName = strings.Replace(productName, "<!-- -->", "", -1)
			res.ProductName = strings.Replace(res.ProductName, "&amp;", "&", -1)
			res.ProductName = strings.Replace(res.ProductName, "&#39;", "'", -1)
		}

		intPrice, errInt := body.DOM.Find(carrefourIntPriceElementQuerySelector.String()).Html()
		decimalPrice, errDecimal := body.DOM.Find(carrefourDecimalPriceElementQuerySelector.String()).Html()

		if errInt == nil && errDecimal == nil {
			priceString := intPrice + "." + decimalPrice
			price, err := strconv.ParseFloat(priceString, 32)

			if err != nil {
				fmt.Printf("Error crawling %s : %s\n", url, err)
				resCh <- crawl.ResultModel{CrawlUrl: ""}
				return
			}

			res.ProductPrice = float32(price)
		}

		imageUrl := body.ChildAttr(carrefourImageElementQuerySelector.String(), "src")

		res.Store = crawler.store
		res.ImageUrl = imageUrl

		resCh <- res
		fmt.Printf("Carrefour from url %s: %v\n", url, res)
	})

	err := collyClient.Visit(url)
	if err != nil {
		fmt.Printf("Error crawling %s : %s\n", url, err)
		resCh <- crawl.ResultModel{CrawlUrl: ""}
		return
	}
}
