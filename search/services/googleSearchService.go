package services

import (
	"fmt"
	"lib/data/dto/crawl"
	"lib/data/dto/store"
	"lib/functional"
	"lib/network/http"
	"net/url"
	"os"
	dtoTypes "search/data/dto"
	"search/data/repositories"
)

type GoogleSearchService struct {
	storeRepo *repositories.StoreMetadata
}

func NewGoogleSearchService(storeRepo *repositories.StoreMetadata) *GoogleSearchService {
	return &GoogleSearchService{storeRepo: storeRepo}
}

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
	res, err := http.GetSync[dtoTypes.GoogleSearchDto](requestUrl)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func (searchService GoogleSearchService) SearchCrawlSources(query string) ([]crawl.SourceDto, error) {
	searchResult, err := queryGoogle(query)

	if err != nil {
		return nil, err
	}

	crawlSources := functional.Map[dtoTypes.GoogleSearchItemDto, crawl.SourceDto](
		searchResult.Items,
		func(item dtoTypes.GoogleSearchItemDto) crawl.SourceDto {
			parsedUrl, err := url.Parse(item.Link)
			if err != nil {
				return crawl.SourceDto{
					Url: "",
					Store: store.MetadataDto{
						StoreId: 0,
					},
				}
			}
			return crawl.SourceDto{
				Url:   item.Link,
				Store: searchService.storeRepo.GetForUrl(parsedUrl.Host).ToDto(),
			}
		})

	filteredCrawlSources := functional.Filter[crawl.SourceDto](
		crawlSources,
		func(source crawl.SourceDto) bool {
			return source.Store.StoreId != 0
		})

	return filteredCrawlSources, nil
}
