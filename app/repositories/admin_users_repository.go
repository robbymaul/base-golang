package repositories

import (
	"context"
	"fmt"
	"paymentserviceklink/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertAdminUserRepository(ctx context.Context, adminUsers *models.AdminUsers) (*models.AdminUsers, error) {
	var (
		db = rc.db
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)

	err := db.Create(&adminUsers).Error
	if err != nil {
		return nil, err
	}

	return adminUsers, nil
}

func (rc *RepositoryContext) UpdateSequenceSetValRepository(ctx context.Context) error {
	var (
		db = rc.db
	)

	// Corrected SQL query syntax
	err := db.Exec("SELECT setval('admin_users_id_seq', (SELECT MAX(id) FROM admin_roles))").Error
	if err != nil {
		return fmt.Errorf("failed to update sequence: %w", err)
	}

	return nil
}

func (rc *RepositoryContext) GetAdminUserByIdRepository(ctx context.Context, id string) (adminUser *models.AdminUsers, err error) {
	var (
		db = rc.db
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)

	db = db.Where("uuid = ?", id)

	db = db.Preload("AdminRole")

	err = db.First(&adminUser).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get admin user by id")
		return nil, err
	}

	return adminUser, err
}

func (rc *RepositoryContext) GetAdminUserByCLIRepository(ctx context.Context, id int64) (adminUser *models.AdminUsers, err error) {
	var (
		db = rc.db
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)

	db = db.Where("id = ?", id)

	db = db.Preload("AdminRole")

	err = db.First(&adminUser).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get admin user by id")
		return nil, err
	}

	return adminUser, err
}

func (rc *RepositoryContext) GetAdminUserRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.AdminUsers, error) {
	var (
		db        = rc.db
		adminUser *models.AdminUsers
		err       error
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&adminUser).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get admin user")
		return nil, err
	}

	return adminUser, nil
}

func (rc *RepositoryContext) UpdateAdminUserRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (err error) {
	var (
		db = rc.db
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)
	db = fn(db)

	return db.Error
}

func (rc *RepositoryContext) FindAdminUserRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.AdminUsers, error) {
	var (
		db         = rc.db
		adminUsers []*models.AdminUsers
		err        error
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&adminUsers).Error
	if err != nil {
		log.Err(err).Msg("error query find admin user")
		return nil, err
	}

	return adminUsers, nil
}

func (rc *RepositoryContext) GetCountListAdminUserRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int, error) {
	var (
		db         = rc.db
		totalCount int64
		err        error
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&totalCount).Error
	if err != nil {
		log.Error().Err(err).Msg("error query select count(*) admin user")
		return 0, err
	}

	return int(totalCount), nil
}

func (rc *RepositoryContext) DeleteAdminUserRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	db = db.Table("admin_users").Model(&models.AdminUsers{}).WithContext(ctx)
	db = fn(db)

	return db.Error
}
