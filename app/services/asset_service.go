package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	"strings"
)

type AssetService struct {
	service     *Service
	serviceName string
}

func NewAssetService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *AssetService {
	return &AssetService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "asset_service",
	}
}

func (s *AssetService) UploadFileService(file multipart.File, header *multipart.FileHeader, assetType string) (*web.UploadFileResponse, error) {
	payload := web.UploadFileRequest{
		AssetType: enums.AssetType(assetType),
		FileName:  header.Filename,
	}

	ext := filepath.Ext(header.Filename)
	log.Debug().Interface("ext", ext).Msg("file extension")

	// validation extension file
	err := s.isExtAllowed(payload.AssetType, ext)
	if err != nil {
		log.Error().Err(err).Msg("is extension allowed error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	// set base file name uuid
	baseFileName := uuid.New().String()
	payload.FileName = baseFileName + ext

	// get asset path
	assetPath, err := s.GetAssetPath(payload.AssetType, payload.FileName)
	if err != nil {
		log.Error().Err(err).Msg("get asset path error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	_, err = s.service.repository.UploadFileRepository(s.service.ctx, assetPath, file, header.Size)
	if err != nil {
		log.Error().Err(err).Msg("upload file error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return &web.UploadFileResponse{
		FileName: payload.FileName,
		FileUrl:  s.urlFile(assetPath),
	}, nil
}

func (s *AssetService) isExtAllowed(assetType enums.AssetType, ext string) error {
	var allowedMap map[string]bool

	switch assetType {
	case enums.ASSET_LOGO:
		allowedMap = map[string]bool{
			"png":  true,
			"jpg":  true,
			"jpeg": true,
		}
	default:
		return fmt.Errorf("image %v extension is not allowed", assetType)
	}

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	log.Debug().Interface("ext", ext).Msg("ext")
	if allowed, ok := allowedMap[ext]; !ok || !allowed {
		return fmt.Errorf("image %v extension is not allowed", assetType)
	}

	return nil
}

func (s *AssetService) GetAssetPath(assetType enums.AssetType, name string) (string, error) {
	assetDir := ""

	switch assetType {
	case enums.ASSET_LOGO:
		assetDir = s.service.config.S3AssetLogo
	default:
		return "", fmt.Errorf("unknown asset type and asset path")
	}

	return path.Join(assetDir, name), nil
}

func (s *AssetService) urlFile(assetPath string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.service.config.S3BucketName, s.service.config.S3Region, assetPath)
}
