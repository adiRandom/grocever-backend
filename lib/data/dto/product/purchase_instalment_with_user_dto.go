package product

type PurchaseInstalmentWithUserDto struct {
	PurchaseInstalmentDto
	UserId int `json:"userId"`
}
