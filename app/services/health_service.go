package services

import (
	"context"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	"time"
)

type HealthService struct {
	service   *Service
	StartTime time.Time
	config    *config.Config
}

func NewHealthService(ctx context.Context, startTime time.Time, r *repositories.RepositoryContext, config *config.Config) *HealthService {
	return &HealthService{
		service:   NewService(ctx, r, config),
		StartTime: startTime,
		config:    config,
	}
}

func (s *HealthService) GetHealthService() *web.HealthResponseWeb {
	return &web.HealthResponseWeb{
		AppName:    s.config.Application,
		Uptime:     time.Since(s.StartTime).String(),
		AppVersion: "N/A",
		Resource: web.ResourcesWeb{
			Database: "OK",
		},
	}
}
