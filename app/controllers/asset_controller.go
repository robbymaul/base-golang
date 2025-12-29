package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
)

func (c *Controller) AdminUploadFileImageController(ctx *gin.Context) {
	context := "AssetController.AdminUploadFile"

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("upload file error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("header", header).Interface("file", file).Interface("context", context).Msg("upload file success")

	defer func() {
		if errC := file.Close(); errC != nil {
			log.Error().Err(errC).Interface("context", context).Msg("close file error")
			helpers.ErrorResponse(ctx, helpers.NewErrorTrace(errC, context).WithStatusCode(http.StatusInternalServerError))
			return
		}
	}()

	assetType := ctx.PostForm("assetType")
	if assetType == "" {
		log.Error().Interface("context", context).Msg("asset type is empty")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.UploadFileImageValidation(file, header, assetType)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("upload file validation error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	assetService := services.NewAssetService(ctx, c.repo, c.cfg)

	response, err := assetService.UploadFileService(file, header, assetType)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("upload file service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})

}
