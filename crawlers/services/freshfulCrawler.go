package services

import (
	types "crawlers/data/dto"
	"fmt"
	"lib/data/models"
	"lib/data/models/crawl"
	"lib/helpers"
	"lib/network/http"
	url2 "net/url"
	"strings"
)

type FreshfulCrawler struct {
	store models.StoreMetadata
}

const freshfulApiUrl = "https://www.freshful.ro/api/v2/shop/product-by-slug/%s"

func getFreshfulProductUrl(url string) (*string, error) {
	parsedUrl, err := url2.Parse(url)
	if err != nil {
		return nil, helpers.Error{Msg: fmt.Sprintf(cannotParseUrlError, "freshful", url), Reason: notUrlErrorReason}
	}

	path := parsedUrl.Path
	segments := strings.Split(path, "/")

	slug, err := helpers.SafeGet(segments, 2)

	if err != nil {
		return nil, helpers.Error{Msg: fmt.Sprintf(cannotParseUrlError, "freshful", url), Reason: notEnoughSegmentsErrorReason}
	}

	correctUrl := fmt.Sprintf(freshfulApiUrl, *slug)

	return &correctUrl, nil
}

func (c FreshfulCrawler) ScrapeProductPage(url string, resCh chan crawl.ResultModel) {
	correctUrl, err := getFreshfulProductUrl(url)
	if err != nil {
		return
	}

	apiRes, err := http.GetSync[types.FreshfulDto](*correctUrl)
	if err != nil {
		return
	}

	res := crawl.ResultModel{CrawlUrl: url}
	res.ProductName = apiRes.Name
	res.ProductPrice = apiRes.Price
	res.Store = c.store

	resCh <- res
}
