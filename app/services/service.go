package services

import (
	"application/app/repositories"
	"application/config"
	"context"
)

type Service struct {
	ctx        context.Context
	config     *config.Config
	repository *repositories.RepositoryContext
}

func NewService(ctx context.Context, r *repositories.RepositoryContext, cfg *config.Config) *Service {

	return &Service{ctx: ctx, config: cfg, repository: r}
}
