package crawlers

import (
	types "dealScraper/crawlers/data/dto"
	"dealScraper/crawlers/models"
	"dealScraper/lib/data/constants"
	"dealScraper/lib/helpers"
	"dealScraper/lib/network"
	"fmt"
	url2 "net/url"
	"strings"
)

type FreshfulCrawler struct {
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

func (c FreshfulCrawler) ScrapeProductPage(url string, resCh chan models.CrawlerResult) {
	correctUrl, err := getFreshfulProductUrl(url)
	if err != nil {
		return
	}

	apiRes, err := network.GetSync[types.FreshfulDto](*correctUrl)
	if err != nil {
		return
	}

	res := models.CrawlerResult{CrawlUrl: url}
	res.ProductName = apiRes.Name
	res.ProductPrice = apiRes.Price
	res.StoreId = constants.FreshfulStoreId

	resCh <- res
}
