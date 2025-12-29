package controllers

import (
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/pagination"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	repo       *repositories.RepositoryContext
	cfg        *config.Config
	startTime  time.Time
	appVersion string
	context    string
}

func NewController(startTime time.Time, appVersion string, cfg *config.Config, repo *repositories.RepositoryContext) *Controller {
	return &Controller{
		cfg:        cfg,
		repo:       repo,
		startTime:  startTime,
		appVersion: appVersion,
	}
}

func (c *Controller) Pagination(ctx *gin.Context) (*pagination.Pages, error) {
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	order := ctx.DefaultQuery("order", "asc")
	//search := ctx.DefaultQuery("search", "")
	isDelete := ctx.DefaultQuery("is_deleted", "false")
	filter := ctx.Query("filter")
	joinOperator := ctx.DefaultQuery("join_operator", "and")

	pages, err := pagination.New(page, perPage, 0, order, isDelete, filter, joinOperator)
	if err != nil {
		return nil, err
	}
	log.Debug().Interface("pages", pages).Interface("context", c.context).Msg("pages pagination init")

	return pages, nil
}
