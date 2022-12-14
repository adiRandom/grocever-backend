package services

import (
	"fmt"
	"lib/functional"
	"lib/helpers"
	"ocr/api/store"
	"ocr/models"
	"ocr/utils"
	"regexp"
	"strconv"
	"strings"
)

var priceLineRegex = regexp.MustCompile(`((\d+)|(\d+(\.|,)\d{1,3})) * ((BUC)|(KG))\.? *(X|x) *\d+(\.|,)\d{2}`)
var qtyRegex = regexp.MustCompile(`^((\d+)|(\d+\.\d{1,3}))`)
var unitRegex = regexp.MustCompile(`(BUC)|(KG)`)
var unitPriceRegex = regexp.MustCompile(`\d+(\.|,)\d{2}$`)

type ParseService struct {
	storeApi *store.Client
}

var parseService *ParseService = nil

func GetParseService() ParseService {
	if parseService == nil {
		parseService = &ParseService{
			storeApi: store.GetClient(),
		}
	}
	return *parseService
}

func (s *ParseService) GetOcrProducts(ocrText string) ([]models.OcrProduct, error) {
	storeName, err := s.getStore(ocrText)
	if err != nil {
		return nil, err
	}

	storeMetadata, err := s.storeApi.GetStoreMetadataForName(storeName)
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(ocrText, "\n")
	// Remove the header
	tokens = tokens[storeMetadata.OcrHeaderLines:]
	productAndPrice := s.zipProductAndPrice(tokens)

	products, err := s.getOcrProductsFromPairs(productAndPrice, storeMetadata.StoreId)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ParseService) getOcrProductsFromPairs(
	productAndPrice []helpers.Pair[string, string],
	storeId int,
) ([]models.OcrProduct, error) {
	products := make([]models.OcrProduct, len(productAndPrice))
	for i, pair := range productAndPrice {
		product := pair.First
		priceLine := pair.Second

		qty, err := s.getQty(priceLine)
		if err != nil {
			return nil, helpers.Error{
				Msg:    fmt.Sprintf("Could not parse qty for %s", product),
				Reason: err.Error(),
			}
		}
		unit := s.getUnit(priceLine)
		unitPrice, err := s.getUnitPrice(priceLine)
		if err != nil {
			return nil, helpers.Error{
				Msg:    fmt.Sprintf("Could not parse unit price for %s", product),
				Reason: err.Error(),
			}
		}

		products[i] = models.NewOcrProduct(
			product,                              // name
			unit,                                 // unitName
			float32(utils.TruncateFloat(qty, 3)), // qty
			float32(utils.TruncateFloat(unitPrice, 3)), // unitPrice,
			storeId, // storeId
		)
	}
	return products, nil
}

func (s *ParseService) getStore(ocrText string) (string, error) {
	allCapsOcrText := strings.ToUpper(ocrText)
	storeNames := s.storeApi.GetAllStoreNames()

	storeNameRegexStr := ""
	for i, name := range storeNames {
		storeNameRegexStr += "(" + name + ")"
		if i < len(storeNames)-1 {
			storeNameRegexStr += "|"
		}
	}

	// Get store name from ocrText
	regex, err := regexp.Compile(storeNameRegexStr)
	if err != nil {
		return "", err
	}

	match := regex.FindString(allCapsOcrText)
	if match == "" {
		return "", helpers.Error{Msg: "No store name found in ocr text"}
	}

	return match, nil
}

/*
* Returns an array of pairs where the first element of the pair is the product name
* and the second element is the price line
 */
func (s *ParseService) zipProductAndPrice(tokens []string) []helpers.Pair[string, string] {
	upperCaseTokens := functional.Map(tokens, func(token string) string {
		return strings.ToUpper(token)
	})
	isFirstLinePrice := priceLineRegex.MatchString(upperCaseTokens[0])
	// Count how many tokens are price lines
	productCount := 0
	for _, token := range tokens {
		if priceLineRegex.MatchString(token) {
			productCount++
		}
	}

	pairs := make([]helpers.Pair[string, string], productCount)
	pairsIndex := 0
	i := 0
	for pairsIndex < productCount {
		if isFirstLinePrice {
			pairs[pairsIndex] = helpers.Pair[string, string]{tokens[i+1], tokens[i]}
		} else {
			pairs[pairsIndex] = helpers.Pair[string, string]{tokens[i], tokens[i+1]}
		}
		i += 2
		pairsIndex++
	}

	return pairs
}

func (s *ParseService) getQty(priceLine string) (float64, error) {
	match := qtyRegex.FindString(priceLine)
	match = strings.ReplaceAll(match, ",", ".")
	return strconv.ParseFloat(match, 32)
}

func (s *ParseService) getUnit(priceLine string) string {
	return unitRegex.FindString(priceLine)
}

func (s *ParseService) getUnitPrice(priceLine string) (float64, error) {
	trimmedLine := strings.Trim(priceLine, " ")
	match := unitPriceRegex.FindString(trimmedLine)
	match = strings.ReplaceAll(match, ",", ".")
	return strconv.ParseFloat(match, 32)
}
