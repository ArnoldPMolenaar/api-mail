package requests

type SendMailAttachment struct {
	FileName string `json:"fileName" validate:"required"`
	FileType string `json:"fileType" validate:"required"`
	FileSize int64  `json:"fileSize" validate:"required"`
	FileData []byte `json:"fileData" validate:"required"`
}
