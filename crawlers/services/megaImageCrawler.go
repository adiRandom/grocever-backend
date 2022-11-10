package services

import (
	types "crawlers/data/dto"
	"fmt"
	"lib/data/constants"
	"lib/data/models"
	"lib/helpers"
	"lib/network/http"
	url2 "net/url"
	"strings"
)

const megaImageApiUrl = "https://api.mega-image.ro/?operationName=ProductDetails&variables={\"productCode\":\"%s\",\"lang\":\"ro\"}&extensions={\"persistedQuery\":{\"version\":1,\"sha256Hash\":\"c734fc7b27b17d674c66d9ae0c70caf29dbdf2667b0300e50f89fc444418d59b\"}}"

type MegaImageCrawler struct {
}

func getMegaImageProductUrl(url string) (*string, error) {
	parsedUrl, err := url2.Parse(url)
	if err != nil {
		return nil, helpers.Error{Msg: fmt.Sprintf(cannotParseUrlError, "Mega Image", url), Reason: notUrlErrorReason}
	}

	path := parsedUrl.Path
	segments := strings.Split(path, "/")

	productId, err := helpers.SafeGet(segments, len(segments)-1)

	if err != nil {
		return nil, helpers.Error{Msg: fmt.Sprintf(cannotParseUrlError, "Mega Image", url), Reason: notEnoughSegmentsErrorReason}
	}

	correctUrl := fmt.Sprintf(megaImageApiUrl, *productId)

	return &correctUrl, nil
}

func (crawler MegaImageCrawler) ScrapeProductPage(url string, resCh chan models.CrawlerResult) {
	correctUrl, err := getMegaImageProductUrl(url)
	if err != nil {
		return
	}

	apiRes, err := http.GetSync[types.MegaImageDto](*correctUrl)
	if err != nil {
		return
	}

	res := models.CrawlerResult{CrawlUrl: url}
	res.ProductName = apiRes.Data.ProductDetails.Name
	res.ProductPrice = apiRes.Data.ProductDetails.Price.Value
	res.StoreId = constants.MegaImageStoreId

	resCh <- res
}
