package dto

type TextSimilarityDto struct {
	Text1 string `json:"text_1"`
	Text2 string `json:"text_2"`
}

type TextSimilarityResultDto struct {
	Similarity float64 `json:"similarity"`
}
