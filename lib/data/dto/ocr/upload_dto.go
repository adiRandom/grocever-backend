package ocr

import "mime/multipart"

type UploadDto struct {
	Bytes  []byte
	Size   int64
	UserId int
}

type UploadImageRequest struct {
	Image  *multipart.File
	UserId int
}

func NewUploadImageRequest(image *multipart.File, userId int) *UploadImageRequest {
	return &UploadImageRequest{Image: image, UserId: userId}
}

const UploadImageParam = "image"
const UploadImageUserIdParam = "userId"
