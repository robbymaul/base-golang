package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	pkgjwt "paymentserviceklink/pkg/jwt"
	"paymentserviceklink/pkg/middleware"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func GetAdminSessionAuth(s *Service) (*models.AdminUsers, error) {
	session := MustGetUserSession(s.ctx)
	if session == nil {
		return nil, helpers.NewErrorTrace(errors.New("admin user session forbidden"), "").WithStatusCode(http.StatusForbidden)
	}
	adminUser, err := s.repository.GetAdminUserRepository(s.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_users.email," +
			"admin_users.full_name," +
			"admin_users.phone," +
			"admin_users.avatar_url," +
			"admin_users.is_active," +
			"admin_users.is_verified," +
			"admin_users.uuid," +
			"admin_users.role_id")
		db = db.Where("admin_users.username = ?", session.Sub)
		db = db.Preload("AdminRole", func(db *gorm.DB) *gorm.DB {
			db = db.Select("admin_roles.id, admin_roles.code, admin_roles.is_active, admin_roles.permissions")
			return db
		})

		return db
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg(fmt.Sprintf("Admin is not found. username = %s", err.Error()))
			return nil, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusNotFound)
		}

		log.Error().Err(err).Msg(fmt.Sprintf("failed to get admin user by email or username: %s", err.Error()))
		return nil, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusInternalServerError)
	}

	return adminUser, err
}

func MustGetUserSession(ctx context.Context) *pkgjwt.JwtResponse {
	val, ok := ctx.Value(middleware.BearerToken).(*pkgjwt.JwtResponse)
	if !ok {
		return nil
	}

	return val
}

func GetClientAuth(s *Service) (*models.Platforms, error) {
	session := MustGetClientAuth(s.ctx)
	if session == nil {
		return nil, helpers.NewErrorTrace(errors.New("client session forbidden"), "").WithStatusCode(http.StatusForbidden)
	}

	return session, nil
}

func MustGetClientAuth(ctx context.Context) *models.Platforms {
	val, ok := ctx.Value(middleware.Client).(*models.Platforms)
	if !ok {
		return nil
	}

	return val
}
