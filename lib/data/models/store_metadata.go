package models

import "lib/data/dto"

type StoreMetadata struct {
	StoreId        int    `json:"storeId"`
	Name           string `json:"name"`
	OcrHeaderLines int    `json:"ocrHeaderLines"`
	Url            string `json:"url"`
}

func NewStoreMetadataFromDto(dto dto.StoreMetadata) StoreMetadata {
	return StoreMetadata{
		StoreId:        dto.StoreId,
		Name:           dto.Name,
		OcrHeaderLines: dto.OcrHeaderLines,
		Url:            dto.Url,
	}
}

func (s *StoreMetadata) ToDto() dto.StoreMetadata {
	return dto.StoreMetadata{
		StoreId:        s.StoreId,
		Name:           s.Name,
		OcrHeaderLines: s.OcrHeaderLines,
		Url:            s.Url,
	}
}