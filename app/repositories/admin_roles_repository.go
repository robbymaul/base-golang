package repositories

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"paymentserviceklink/app/models"
)

func (rc *RepositoryContext) InsertAdminRoleRepository(ctx context.Context, adminRole *models.AdminRoles) error {
	var (
		db = rc.db
	)

	db = db.Table("admin_roles").Model(&models.AdminRoles{}).WithContext(ctx)

	return db.Create(&adminRole).Error
}

func (rc *RepositoryContext) FindAdminRoleRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.AdminRoles, error) {
	var (
		db         = rc.db
		adminRoles []*models.AdminRoles
		err        error
	)

	db = db.Table("admin_roles").Model(&models.AdminRoles{}).WithContext(ctx)
	db = fn(db)

	db = db.Order("id asc")

	err = db.Find(&adminRoles).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get list admin role")
		return nil, err
	}

	return adminRoles, nil
}

func (rc *RepositoryContext) GetTotalCountAdminRoleRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db         = rc.db
		totalCount int64
		err        error
	)

	db = db.Table("admin_roles").Model(&models.AdminRoles{}).WithContext(ctx)
	db = fn(db)

	db = db.Order("id asc")

	err = db.Count(&totalCount).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get count admin role")
		return 0, err
	}

	return totalCount, nil
}

func (rc *RepositoryContext) GetAdminRoleRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.AdminRoles, error) {
	var (
		db        = rc.db
		adminRole *models.AdminRoles
		err       error
	)

	db = db.Table("admin_roles").Model(&models.AdminRoles{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&adminRole).Error
	if err != nil {
		log.Err(err).Msg("error query get admin role")
		return nil, err
	}

	return adminRole, nil
}

func (rc *RepositoryContext) UpdateAdminRoleRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)
	db = db.Table("admin_roles").Model(&models.AdminRoles{}).WithContext(ctx)
	db = fn(db)

	return db.Error
}
