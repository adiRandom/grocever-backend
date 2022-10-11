package utils

import (
	"dealScraper/crawlers/data/constants"
	"net/url"
)

// GetStoreIdForDomain returns the store id for a given domain or 0 if invalid url domain
func GetStoreIdForDomain(url url.URL) int {
	switch url.Host {
	case "www.cora.ro":
		return constants.CoraStoreId
	case "www.freshful.ro":
		return constants.FreshfulStoreId
	case "www.mega-image.ro":
		return constants.MegaImageStoreId
	case "www.auchan.ro":
		return constants.AuchanStoreId
	default:
		return 0
	}
}
