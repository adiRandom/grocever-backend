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

const sezamoContentElementQuerySelector = utils.CssSelector("[data-test=product-detail-upper-section]")
const sezamoTitleElementQuerySelector = utils.CssSelector("[data-test=product-detail-upper-section] div:nth-child(3) h2 a")
const sezamoPriceElementQuerySelector = utils.CssSelector("[data-test=product-price]")
const sezamoImageElementQuerySelector = utils.CssSelector("[data-gtm-item=product-image] div picture img")

type SezamoCrawler struct {
	store models.StoreMetadata
}

func (crawler SezamoCrawler) ScrapeProductPage(url string, resCh chan crawl.ResultModel) {
	collyClient := colly.NewCollector()

	collyClient.OnHTML("body", func(body *colly.HTMLElement) {
		if body.DOM.Find(sezamoContentElementQuerySelector.String()).Length() == 0 {
			// Not url to product page
			resCh <- crawl.ResultModel{CrawlUrl: ""}
			return
		}

		res := crawl.ResultModel{CrawlUrl: url}
		productName, err := body.DOM.Find(sezamoTitleElementQuerySelector.String()).Html()

		if err == nil {
			res.ProductName = strings.Replace(productName, "<!-- -->", "", -1)
		}

		price, err := body.DOM.Find(sezamoPriceElementQuerySelector.String()).Html()
		if err != nil {
			fmt.Printf("Error crawling %s : %s\n", url, err)
			resCh <- crawl.ResultModel{CrawlUrl: ""}
			return
		}

		price = strings.Split(price, "<!-- -->")[1]
		price = strings.Replace(price, "<!-- -->", "", -1)
		// There is a wierd character at index 4 that is not a space
		// Substring to remove it
		price = strings.Replace(price, "lei", "", -1)[:4]
		price = strings.Replace(price, ",", ".", -1)

		parsedPrice, err := strconv.ParseFloat(price, 32)
		res.ProductPrice = float32(parsedPrice)

		imageUrl := body.ChildAttr(sezamoImageElementQuerySelector.String(), "src")

		res.Store = crawler.store
		res.ImageUrl = imageUrl

		resCh <- res
		fmt.Printf("sezamoCrawler from url %s: %v\n", url, res)
	})

	err := collyClient.Visit(url)
	if err != nil {
		fmt.Printf("Error crawling %s : %s\n", url, err)
		resCh <- crawl.ResultModel{CrawlUrl: ""}
		return
	}
}
