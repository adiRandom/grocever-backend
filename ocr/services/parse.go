package services

import (
	libModels "lib/data/models"
	"lib/data/models/product"
	"lib/helpers"
	"ocr/api/store"
	"ocr/utils"
	"regexp"
	"strconv"
	"strings"
)

var priceLineRegex = regexp.MustCompile(`((\d+)|(\d+(\.|,)\d{1,3})) * ((BUC)|(KG))\.? *(X|x) *\d+(\.|,)\d{2}`)
var qtyRegex = regexp.MustCompile(`^((\d+)|(\d+\.\d{1,3}))`)
var unitRegex = regexp.MustCompile(`(BUC)|(KG)`)
var unitPriceRegex = regexp.MustCompile(`\d+(\.|,)\d{2}$`)
var unitAndUnitPriceRegex = regexp.MustCompile(`((BUC)|(KG))\.? *(X|x) *\d+(\.|,)\d{2}`)
var unitAndUnitPriceBeginningRegex = regexp.MustCompile(`^((BUC)|(KG))\.? *(X|x) *\d+(\.|,)\d{2}`)
var qtyBeginningRegex = regexp.MustCompile(`^((\d+)|(\d+\.\d{1,3}))`)

var skipLinesMarkers = []string{
	"C.I.F",
	"CIF",
	"COD IDENTIFICARE FISCALA",
	"LEI",
	"RON",
	"BON FISCAL",
}

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

func (s *ParseService) GetOcrProducts(ocrText string, userId int) ([]product.PurchaseInstalmentModel, error) {
	storeName, err := s.getStore(ocrText)
	if err != nil {
		return nil, err
	}

	storeMetadataDto, err := s.storeApi.GetStoreMetadataForName(storeName)
	storeMetadata := libModels.NewStoreMetadataFromDto(storeMetadataDto)
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(ocrText, "\n")
	// Remove the header
	tokens = s.getRelevantTokens(tokens)
	productAndPrice := s.zipProductAndPrice(tokens)

	products := s.getOcrProductsFromPairs(productAndPrice, storeMetadata, userId)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ParseService) getRelevantTokens(tokens []string) []string {
	// Get the first line that matches the price line regex
	for index, token := range tokens {
		if priceLineRegex.MatchString(token) {
			upperPrevToken := strings.ToUpper(tokens[index-1])
			for _, marker := range skipLinesMarkers {
				if strings.Contains(upperPrevToken, marker) {
					// The line above the first price line is a skipped line
					return tokens[index:]
				}
			}
			// The line above the first price line is not a skipped line
			return tokens[index-1:]
		}
	}

	return []string{}
}

func (s *ParseService) getOcrProductsFromPairs(
	productAndPrice []helpers.Pair[string, string],
	store libModels.StoreMetadata,
	userId int,
) []product.PurchaseInstalmentModel {
	products := make([]product.PurchaseInstalmentModel, len(productAndPrice))
	for i, pair := range productAndPrice {
		ocrProductName := pair.First
		priceLine := pair.Second

		qty, err := s.getQty(priceLine)
		if err != nil {
			continue
		}
		unit := s.getUnit(priceLine)
		unitPrice, err := s.getUnitPrice(priceLine)
		if err != nil {
			continue
		}

		price := float32(utils.TruncateFloat(qty, 3)) * float32(utils.TruncateFloat(unitPrice, 3))

		ocrProduct := product.NewOcrProductModel(ocrProductName, nil, nil, nil)

		products[i] = *product.NewPurchaseInstalmentModel(
			-1,                                   // id
			float32(utils.TruncateFloat(qty, 3)), // qty
			price,                                // price
			userId,                               // userId
			*ocrProduct,                          // ocrProduct
			float32(utils.TruncateFloat(unitPrice, 3)), // unitPrice
			store, // store
			unit,  // unit
		)
	}
	return products
}

func (s *ParseService) getStore(ocrText string) (string, error) {
	allCapsOcrText := strings.ToUpper(ocrText)
	storeNames := s.storeApi.GetAllStoreNames()

	storeNameRegexStr := ""
	for i, name := range storeNames {
		storeNameRegexStr += "(" + strings.ToUpper(name) + ")"
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
	isFirstLinePrice := priceLineRegex.MatchString(strings.ToUpper(tokens[0]))
	reconciledTokens := s.reconcileSplitTokens(tokens, isFirstLinePrice)

	// Count how many tokens are price lines
	productCount := 0
	for _, token := range reconciledTokens {
		if priceLineRegex.MatchString(strings.ToUpper(token)) {
			productCount++
		}
	}

	pairs := make([]helpers.Pair[string, string], productCount)
	pairsIndex := 0
	i := 0
	for pairsIndex < productCount {
		if isFirstLinePrice {
			pairs[pairsIndex] = helpers.Pair[string, string]{reconciledTokens[i+1], reconciledTokens[i]}
		} else {
			pairs[pairsIndex] = helpers.Pair[string, string]{reconciledTokens[i], reconciledTokens[i+1]}
		}
		i += 2
		pairsIndex++
	}

	return pairs
}

// If the price line was split wrong, merge them
// If the product name was split in multiple lines, merge them
func (s *ParseService) reconcileSplitTokens(tokens []string, isFirstLinePriceLine bool) []string {
	// copy the tokens
	newTokens := make([]string, len(tokens))
	copy(newTokens, tokens)

	for i := 0; i < len(newTokens); i++ {
		if (i%2 == 0 && isFirstLinePriceLine) || (i&2 == 1 && !isFirstLinePriceLine) {
			token := newTokens[i]
			qtyMatch := qtyRegex.MatchString(token)
			unitAndUnitPriceMatch := unitAndUnitPriceRegex.MatchString(strings.ToUpper(token))

			if qtyMatch && unitAndUnitPriceMatch {
				continue
			}

			if qtyMatch {
				for j := i + 1; j < len(newTokens); j++ {
					if unitAndUnitPriceBeginningRegex.MatchString(strings.ToUpper(newTokens[j])) {
						newTokens[i] = token + " " + (newTokens)[j]
						newTokens = append((newTokens)[:j], (newTokens)[j+1:]...)
						break
					}
				}
			} else if unitAndUnitPriceMatch {
				for j := i - 1; j >= 0; j-- {
					if qtyBeginningRegex.MatchString((newTokens)[j]) {
						newTokens[i] = (newTokens)[j] + " " + token
						newTokens = append((newTokens)[:j], (newTokens)[j+1:]...)
						break
					}
				}
			} else {
				// This is a product name split in 2 lines
				// Append this to the previous token
				newTokens[i-1] = newTokens[i-1] + " " + token
				newTokens = append(newTokens[:i], newTokens[i+1:]...)
				// We merged the current token with the previous one
				// so we need to decrement i so that we don't skip the next token
				i--
			}
		}
	}

	return newTokens
}

func (s *ParseService) getQty(priceLine string) (float64, error) {
	match := qtyRegex.FindString(priceLine)
	match = strings.ReplaceAll(match, ",", ".")
	return strconv.ParseFloat(match, 32)
}

func (s *ParseService) getUnit(priceLine string) string {
	return unitRegex.FindString(strings.ToUpper(priceLine))
}

func (s *ParseService) getUnitPrice(priceLine string) (float64, error) {
	trimmedLine := strings.Trim(priceLine, " ")
	match := unitPriceRegex.FindString(trimmedLine)
	match = strings.ReplaceAll(match, ",", ".")
	return strconv.ParseFloat(match, 32)
}
