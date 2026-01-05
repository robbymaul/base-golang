package controllers

import (
	"application/app/repositories"
	"application/config"
	"time"
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
