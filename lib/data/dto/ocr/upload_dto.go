package ocr

type UploadDto struct {
	Bytes  []byte
	Size   int64
	UserId uint
}
