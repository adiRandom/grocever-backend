package crawlers

import (
	"dealScraper/crawlers/data/constants"
	"dealScraper/crawlers/models"
	"dealScraper/lib/helpers"
	"dealScraper/lib/network"
	"fmt"
	url2 "net/url"
	"strings"
)

type FreshfulCrawler struct {
}

type feshfulDto struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

const freshfulApiUrl = "https://www.freshful.ro/api/v2/shop/product-by-slug/%s"
const cannotParseUrlError = "cannot get product slug for freshful url: %s"
const notUrlErrorReason = "not a valid url"
const notEnoughSegmentsErrorReason = "not enough segments in url"

func getFreshfulProductUrl(url string) (*string, error) {
	parsedUrl, err := url2.Parse(url)
	if err != nil {
		return nil, helpers.Error{Msg: fmt.Sprintf(cannotParseUrlError, url), Reason: notUrlErrorReason}
	}

	path := parsedUrl.Path
	segments := strings.Split(path, "/")

	slug, err := helpers.SafeGet(segments, 2)

	if err != nil {
		return nil, helpers.Error{Msg: fmt.Sprintf(cannotParseUrlError, url), Reason: notEnoughSegmentsErrorReason}
	}

	correctUrl := fmt.Sprintf(freshfulApiUrl, *slug)

	return &correctUrl, nil
}

func (c FreshfulCrawler) ScrapeProductPage(url string, resCh chan models.CrawlerResult) {
	correctUrl, err := getFreshfulProductUrl(url)
	if err != nil {
		return
	}

	apiRes, err := network.GetSync[feshfulDto](*correctUrl)
	if err != nil {
		return
	}

	res := models.CrawlerResult{}
	res.ProductName = apiRes.Name
	res.ProductPrice = apiRes.Price
	res.StoreId = constants.FreshfulStoreId

	resCh <- res
}
