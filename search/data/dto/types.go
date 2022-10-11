package dto

type GoogleSearchItemDto struct {
	Link string `json:"link"`
}

type GoogleSearchDto struct {
	Items []GoogleSearchItemDto `json:"items"`
}
