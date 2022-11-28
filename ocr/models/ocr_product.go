package models

type OcrProduct struct {
	Name      string
	Price     float32
	UnitName  string
	Qty       float32
	UnitPrice float32
}

func NewOcrProduct(name string, unitName string, qty float32, unitPrice float32) OcrProduct {
	return OcrProduct{
		Name:      name,
		Qty:       qty,
		Price:     qty * unitPrice,
		UnitName:  unitName,
		UnitPrice: unitPrice,
	}
}
