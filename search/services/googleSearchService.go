package services

import (
	"dealScraper/lib/data/dto"
	"dealScraper/lib/functional"
	httpClient "dealScraper/lib/network"
	dtoTypes "dealScraper/search/data/dto"
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

func (searchService GoogleSearchService) SearchCrawlSources(query string) ([]dto.CrawlSourceDto, error) {
	searchResult, err := queryGoogle(query)

	if err != nil {
		return nil, err
	}

	crawlSources := functional.Map[dtoTypes.GoogleSearchItemDto, dto.CrawlSourceDto](
		searchResult.Items,
		func(item dtoTypes.GoogleSearchItemDto) dto.CrawlSourceDto {
			parsedUrl, err := url.Parse(item.Link)
			if err != nil {
				return dto.CrawlSourceDto{
					Url:     "",
					StoreId: 0,
				}
			}
			return dto.CrawlSourceDto{
				Url:     item.Link,
				StoreId: urlUtils.GetStoreIdForDomain(*parsedUrl),
			}
		})

	filteredCrawlSources := functional.Filter[dto.CrawlSourceDto](
		crawlSources,
		func(source dto.CrawlSourceDto) bool {
			return source.StoreId != 0
		})

	return filteredCrawlSources, nil
}
