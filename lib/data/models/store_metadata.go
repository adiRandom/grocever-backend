package models

import (
	"lib/data/dto/store"
)

type StoreMetadata struct {
	StoreId int    `json:"storeId"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

func NewStoreMetadataFromDto(dto store.MetadataDto) StoreMetadata {
	return StoreMetadata{
		StoreId: dto.StoreId,
		Name:    dto.Name,
		Url:     dto.Url,
	}
}

func (s *StoreMetadata) ToDto() store.MetadataDto {
	return store.MetadataDto{
		StoreId: s.StoreId,
		Name:    s.Name,
		Url:     s.Url,
	}
}
