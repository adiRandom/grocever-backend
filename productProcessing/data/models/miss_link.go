package models

import (
	productDto "lib/data/dto/product"
	"productProcessing/data/database/entities"
)

type MissLink struct {
	Id             uint
	ProductId      int
	OcrProductName string
	UserId         int
}

func (model *MissLink) ToEntity() *entities.MissLink {
	entity := entities.MissLink{
		ProductIdFk:      model.ProductId,
		OcrProductNameFk: model.OcrProductName,
		UserId:           model.UserId,
	}

	if model.Id != 0 {
		entity.ID = model.Id
	}

	return &entity
}

func (model *MissLink) ToDto() *productDto.ReportDto {
	return &productDto.ReportDto{
		ProductId:      model.ProductId,
		OcrProductName: model.OcrProductName,
	}
}

func NewMissLinkModelFromEntity(entity *entities.MissLink) *MissLink {
	return &MissLink{
		Id:             entity.ID,
		ProductId:      entity.ProductIdFk,
		OcrProductName: entity.OcrProductNameFk,
		UserId:         entity.UserId,
	}
}
