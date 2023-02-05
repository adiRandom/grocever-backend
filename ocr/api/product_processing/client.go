package product_processing

import (
	"github.com/chebyrash/promise"
	"lib/data/dto/product"
	dto "lib/data/dto/product/ocr"
	"lib/helpers"
	"lib/network/http"
	"os"
)

type Client struct {
	baseUrl string
}

var client *Client

func GetClient() *Client {
	if client == nil {
		client = &Client{
			baseUrl: os.Getenv("PRODUCT_PROCESSING_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) OcrProductsExists(ocrNames []string) ([]bool, error) {

	if ocrNames == nil {
		return nil, helpers.Error{Msg: "ocrNames cannot be nil"}
	}

	res, err := http.UnwrapHttpResponse(http.PostSync[http.Response[dto.ProductExistsResponse]](
		s.baseUrl+"/product/ocr/exists",
		dto.ProductExists{
			OcrNames: ocrNames,
		},
	))

	if err != nil {
		return nil, err
	}

	return res.Exists, nil
}

func (s *Client) CreatePurchaseInstalment(instalment product.CretePurchaseInstalmentDto) (*product.PurchaseInstalmentDto, error) {
	purchaseInstalment, err := http.UnwrapHttpResponse(http.PostSync[http.Response[product.PurchaseInstalmentDto]](
		s.baseUrl+"/product/ocr/instalment",
		instalment,
	))

	if err != nil {
		return nil, err
	}

	return purchaseInstalment, nil
}

func (s *Client) CreatePurchaseInstalmentsAsync(
	instalments product.CreatePurchaseInstalmentListDto,
) *promise.Promise[[]product.PurchaseInstalmentDto] {
	return http.UnwrapHttpAsyncResponse(http.PostAsync[http.Response[[]product.PurchaseInstalmentDto]](
		s.baseUrl+"/product/ocr/instalment/list",
		instalments,
	))
}
