package validate

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
)

func UploadFileImageValidation(file multipart.File, header *multipart.FileHeader, assetType string) error {
	const (
		maxFileSize = 0.5 * 1024 * 1024 // 500kb
	)

	var (
		allowedExtension = map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
	)

	if header.Size > maxFileSize {
		return fmt.Errorf("size file image to large, maximum size is %v. your upload size is %v", maxFileSize, header.Size)
	}

	fileExtension := filepath.Ext(header.Filename)
	if currentValue, ok := allowedExtension[fileExtension]; !ok || !currentValue {
		return fmt.Errorf("unsupported file extension %s", fileExtension)
	}

	helpers.IsInList(enums.AssetType(assetType), []enums.AssetType{
		enums.ASSET_LOGO,
	}...)

	return nil
}
