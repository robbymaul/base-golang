package web

import "paymentserviceklink/app/enums"

type UploadFileRequest struct {
	AssetType enums.AssetType `json:"assetType"`
	FileName  string          `json:"fileName"`
}

type UploadFileResponse struct {
	FileName string `json:"fileName"`
	FileUrl  string `json:"fileUrl"`
}

type ImageRequest struct {
	FileName  string               `json:"fileName,omitempty"`
	FileUrl   string               `json:"fileUrl,omitempty"`
	SizeType  enums.ImageSizeType  `json:"sizeType,omitempty"`
	Geometric enums.ImageGeometric `json:"geometric,omitempty"`
}

type ImageResponse struct {
	Id        uint64               `json:"id,omitempty"`
	FileUrl   string               `json:"fileUrl"`
	SizeType  enums.ImageSizeType  `json:"sizeType"`
	Geometric enums.ImageGeometric `json:"geometric"`
}
