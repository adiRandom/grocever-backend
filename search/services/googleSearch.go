package services

import (
	"dealScraper/lib/functional"
	httpClient "dealScraper/lib/network"
	dtoTypes "dealScraper/search/data/dto"
	types "dealScraper/search/data/models"
	urlUtils "dealScraper/search/utils"
	"fmt"
	"net/url"
	"os"
)

type GoogleSearchService struct{}

const googleSearchUrl = "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s"

func queryGoogle(searchTerm string) (*dtoTypes.GoogleSearchDto, error) {
	// Url encode the search term
	encodedSearchTerm := url.QueryEscape(searchTerm)
	// Create the url
	requestUrl := fmt.Sprintf(googleSearchUrl,
		os.Getenv("GOOGLE_SEARCH_API_KEY"),
		os.Getenv("GOOGLE_SEARCH_CX"),
		encodedSearchTerm)

	// Get the response
	res, err := httpClient.GetSync[dtoTypes.GoogleSearchDto](requestUrl)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func (searchService GoogleSearchService) SearchCrawlSources(query string) ([]types.CrawlSource, error) {
	searchResult, err := queryGoogle(query)

	if err != nil {
		return nil, err
	}

	crawlSources := functional.Map[dtoTypes.GoogleSearchItemDto, types.CrawlSource](
		searchResult.Items,
		func(item dtoTypes.GoogleSearchItemDto) types.CrawlSource {
			parsedUrl, err := url.Parse(item.Link)
			if err != nil {
				return types.CrawlSource{
					Url:     "",
					StoreId: 0,
				}
			}
			return types.CrawlSource{
				Url:     item.Link,
				StoreId: urlUtils.GetStoreIdForDomain(*parsedUrl),
			}
		})

	filteredCrawlSources := functional.Filter[types.CrawlSource](
		crawlSources,
		func(source types.CrawlSource) bool {
			return source.StoreId != 0
		})

	return filteredCrawlSources, nil
}
