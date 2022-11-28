package api

import "lib/data/dto"

func GetStoreMetadataForName(name string) (dto.StoreMetadata, error) {
	// TODO: implement
	return dto.StoreMetadata{Name: "MEGA IMAGE", OcrHeaderLines: 4}, nil
}

func GetAllStoreNames() []string {
	// TODO: implement
	return []string{"Auchan", "Carrefour", "Kaufland", "MEGA IMAGE"}
}
