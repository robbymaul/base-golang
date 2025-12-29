package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	pkgjwt "paymentserviceklink/pkg/jwt"
	"paymentserviceklink/pkg/pagination"
	"paymentserviceklink/pkg/util"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type AdminService struct {
	service     *Service
	serviceName string
}

func NewAdminService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *AdminService {
	return &AdminService{service: NewService(ctx, repo, cfg)}
}

func (s *AdminService) AdminLoginService(payload *web.AdminLoginRequest) (*web.AdminSessionResponse, error) {
	s.serviceName = "AdminService.AdminLoginService"

	adminUser, err := s.service.repository.GetAdminUserRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_users.id," +
			"admin_users.username," +
			"admin_users.password_hash," +
			"admin_users.role_id")
		db = db.Where("admin_users.username = ? OR admin_users.email = ?", payload.UsernameOrEmail, payload.UsernameOrEmail)
		db = db.Preload("AdminRole", func(db *gorm.DB) *gorm.DB {
			db = db.Select("admin_roles.id, admin_roles.code, admin_roles.permissions")
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin user by email or username repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("invalid credentials"), s.serviceName).WithStatusCode(http.StatusUnauthorized)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	if adminUser.AdminRole == nil {
		return nil, helpers.NewErrorTrace(fmt.Errorf("user has no assigned role"), s.serviceName).WithStatusCode(http.StatusUnauthorized)
	}

	err = util.ComparePassword(adminUser.Password, payload.Password)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("password and password compare error")
		return nil, helpers.NewErrorTrace(fmt.Errorf("invalid credentials"), s.serviceName).WithStatusCode(http.StatusUnauthorized)
	}

	issuePayload := &pkgjwt.IssueJwtPayload{
		Id:       adminUser.Id,
		Subject:  adminUser.Username,
		Role:     adminUser.AdminRole.Code,
		Lifetime: s.service.config.JwtExpire,
	}

	jwtAdapter := pkgjwt.NewJwtAdapter(s.service.config.JwtIssuer, s.service.config.JwtSecret)

	response, err := jwtAdapter.IssueJwt(issuePayload)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("jwt issue error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return response, nil
}

func (s *AdminService) GetDetailAdmin() (*web.DetailAdminResponse, error) {
	adminUser, err := GetAdminSessionAuth(s.service)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin session auth error")
		return nil, err
	}

	admin := s.mapToDetailAdminResponse(adminUser)

	return admin, nil
}

func (s *AdminService) GetListAdminRoleService(pages *pagination.Pages) (*web.ListResponse, error) {
	var adminRoles []*models.AdminRoles
	var totalCount int64
	var err error

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnAdminRole())
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		adminRoles, err = s.service.repository.FindAdminRoleRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("admin_roles.id," +
				"admin_roles.code," +
				"admin_roles.name," +
				"admin_roles.description," +
				"admin_roles.permissions," +
				"admin_roles.is_active")
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if pages.Unscoped {
				db = db.Unscoped()
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get admin role repository error")
			return err
		}

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetTotalCountAdminRoleRepository(ctx, func(db *gorm.DB) *gorm.DB {
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if pages.Unscoped {
				db = db.Unscoped()
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get total count admin role repository error")
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("async group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = int(totalCount)

	response := make([]*web.AdminRoleResponse, 0, len(adminRoles))

	for _, role := range adminRoles {
		response = append(response, s.mapToAdminRoleResponse(role))
	}

	return &web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *AdminService) CreateAdminRoleService(payload *web.CreateAdminRolesRequest) error {
	adminRole := &models.AdminRoles{
		Code:        payload.Code,
		Name:        payload.Name,
		Description: payload.Description,
		IsActive:    true,
	}

	err := s.service.repository.InsertAdminRoleRepository(s.service.ctx, adminRole)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("insert admin role repository error")
		return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return nil
}

func (s *AdminService) GetDetailAdminRoleService(idInt int) (*web.AdminRoleResponse, error) {
	adminRole, err := s.service.repository.GetAdminRoleRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_roles.id," +
			"admin_roles.code," +
			"admin_roles.name," +
			"admin_roles.description," +
			"admin_roles.permissions," +
			"admin_roles.is_active")
		db = db.Where("admin_roles.id=?", idInt)

		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin role repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.mapToAdminRoleResponse(adminRole), nil
}

func (s *AdminService) UpdateAdminRoleService(payload *web.AdminRoleResponse) (*web.AdminRoleResponse, error) {
	adminRole, err := s.service.repository.GetAdminRoleRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_roles.id," +
			"admin_roles.code," +
			"admin_roles.name," +
			"admin_roles.description," +
			"admin_roles.permissions," +
			"admin_roles.is_active," +
			"admin_roles.updated_at")
		db = db.Where("admin_roles.id=?", payload.Id)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin role repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("admin role", adminRole).Msg("data admin role")

	s.updateAdminRole(&adminRole, payload)
	log.Debug().Interface("admin role update", adminRole).Msg("data admin role update")

	err = s.service.repository.UpdateAdminRoleRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("admin_roles.id = ?", adminRole.Id)
		updateColumn := map[string]interface{}{
			"name":        adminRole.Name,
			"description": adminRole.Description,
			"permissions": adminRole.Permissions,
			"is_active":   adminRole.IsActive,
			"updated_at":  time.Now(),
		}
		db = db.UpdateColumns(updateColumn)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("update admin role repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.mapToAdminRoleResponse(adminRole), nil
}

func (s *AdminService) CreateAdminUserService(session *pkgjwt.JwtResponse, payload *web.CreateAdminUserRequest) (*web.DetailAdminResponse, error) {
	passwordHash, err := util.HashPassword(payload.Password)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("hash password")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	username := fmt.Sprintf("%s%v", strings.ToLower(util.StringWithoutSpace(payload.Username)), time.Now().Year())

	adminUser := &models.AdminUsers{
		UUID:                uuid.New(),
		Username:            username,
		Email:               payload.Email,
		Password:            passwordHash,
		FullName:            payload.FullName,
		Phone:               payload.Phone,
		AvatarUrl:           payload.AvatarUrl,
		RoleId:              payload.RoleId,
		IsActive:            false,
		IsVerified:          false,
		LastLoginAt:         nil,
		LastLoginIp:         "",
		FailedLoginAttempts: 0,
		LockedUntil:         nil,
		PasswordChangedAt:   nil,
		TwoFactorEnabled:    false,
		TwoFactorSecret:     "",
	}

	adminUser, err = s.service.repository.InsertAdminUserRepository(s.service.ctx, adminUser)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("insert admin user repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	adminActivityLog := models.NewAdminActivityLogs(
		session.Id,
		enums.ACTION_ADMIN_ACCTIVITY_CREATE,
		enums.RESOURCE_ADMIN_ACTIVITY_LOG_USER,
		fmt.Sprint(adminUser.UUID),
		fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_USER, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
		enums.NULL_STRING,
		session.Sub,
		payload,
		adminUser,
	)

	err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, nil, adminActivityLog)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
	}

	return s.mapToDetailAdminResponse(adminUser), nil
}

func (s *AdminService) GetListAdminUserService(pages *pagination.Pages) (*web.ListResponse, error) {
	var (
		adminUsers []*models.AdminUsers
		totalCount int
		err        error
	)

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnAdminUser())
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		adminUsers, err = s.service.repository.FindAdminUserRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("admin_users.id," +
				"admin_users.username," +
				"admin_users.email," +
				"admin_users.full_name," +
				"admin_users.phone," +
				"admin_users.avatar_url," +
				"admin_users.role_id," +
				"admin_users.is_active," +
				"admin_users.is_verified," +
				"admin_users.last_login_at," +
				"admin_users.uuid")
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Joins("JOIN admin_roles ON admin_roles.id = admin_users.role_id")
			if pages.Unscoped {
				db = db.Unscoped()
			}
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			db = db.Preload("AdminRole", func(db *gorm.DB) *gorm.DB {
				db = db.Select("admin_roles.id, admin_roles.code, admin_roles.is_active")
				return db
			})
			db = db.Order("id " + pages.Sort)
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get list admin user repository error")
			return err
		}

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetCountListAdminUserRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Joins("JOIN admin_roles ON admin_roles.id = admin_users.role_id")
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if pages.Unscoped {
				db = db.Unscoped()
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get count list admin user repository")
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("error group sync get list admin user service")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = totalCount

	response := make([]*web.DetailAdminResponse, 0, len(adminUsers))

	for _, adminUser := range adminUsers {
		response = append(response, s.mapToDetailAdminResponse(adminUser))
	}

	return &web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *AdminService) GetDetailAdminUserService(id string) (*web.DetailAdminResponse, error) {
	adminUser, err := s.service.repository.GetAdminUserRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_users.id," +
			"admin_users.uuid," +
			"admin_users.full_name," +
			"admin_users.is_active," +
			"admin_users.is_verified," +
			"admin_users.role_id," +
			"admin_users.avatar_url," +
			"admin_users.phone," +
			"admin_users.email," +
			"admin_users.last_login_at",
		)
		db = db.Where("admin_users.uuid=?", id)
		db = db.Preload("AdminRole", func(db *gorm.DB) *gorm.DB {
			db = db.Select("admin_roles.id, admin_roles.code, admin_roles.is_active")
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin user repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.mapToDetailAdminResponse(adminUser), nil
}

func (s *AdminService) UpdateAdminUserService(id string, payload *web.DetailAdminResponse) (*web.DetailAdminResponse, error) {
	adminUser, err := s.service.repository.GetAdminUserRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_users.id," +
			"admin_users.uuid," +
			"admin_users.full_name," +
			"admin_users.is_active," +
			"admin_users.is_verified," +
			"admin_users.role_id," +
			"admin_users.avatar_url," +
			"admin_users.phone," +
			"admin_users.email",
		)
		db = db.Where("admin_users.uuid=?", id)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin user repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	log.Debug().Interface("data admin user", adminUser).Msg("data admin user")

	s.updateAdminUser(&adminUser, payload)

	log.Debug().Interface("data admin user", adminUser).Msg("data admin user")

	err = s.service.repository.UpdateAdminUserRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("admin_users.id = ?", adminUser.Id)
		updateColumn := map[string]interface{}{
			"email":       adminUser.Email,
			"full_name":   adminUser.FullName,
			"phone":       adminUser.Phone,
			"avatar_url":  adminUser.AvatarUrl,
			"role_id":     adminUser.RoleId,
			"is_active":   adminUser.IsActive,
			"is_verified": adminUser.IsVerified,
			"updated_at":  time.Now(),
		}
		db = db.UpdateColumns(updateColumn)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("update admin user repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.mapToDetailAdminResponse(adminUser), nil
}

func (s *AdminService) DeleteAdminUserService(id string) (*web.DetailAdminResponse, error) {
	log.Debug().Str("id", id).Msg("delete admin user service")
	adminUser, err := s.service.repository.GetAdminUserRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("admin_users.id," +
			"admin_users.uuid," +
			"admin_users.full_name," +
			"admin_users.is_active",
		)
		db = db.Where("admin.users.uuid=?", id)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get admin user repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	adminUser.IsActive = false

	err = s.service.repository.DeleteAdminUserRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("admin_users.id = ?", adminUser.Id).Omit("admin_users.id")
		updateColumn := map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
			"deleted_at": time.Now(),
		}
		db = db.UpdateColumns(updateColumn)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("delete admin user repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.mapToDetailAdminResponse(adminUser), nil
}

func (s *AdminService) mapToDetailAdminResponse(adminUser *models.AdminUsers) *web.DetailAdminResponse {
	role := &web.AdminRoleResponse{}

	if adminUser.AdminRole != nil {
		role = s.mapToAdminRoleResponse(adminUser.AdminRole)
	}

	return &web.DetailAdminResponse{
		Id:                adminUser.UUID.String(),
		Username:          adminUser.Username,
		Email:             adminUser.Email,
		FullName:          adminUser.FullName,
		Phone:             adminUser.Phone,
		Avatar:            adminUser.Username,
		RoleId:            adminUser.RoleId,
		IsActive:          adminUser.IsActive,
		IsVerified:        adminUser.IsVerified,
		LastLoginAt:       adminUser.LastLoginAt,
		LastLoginIp:       adminUser.LastLoginIp,
		LockedUntil:       adminUser.LockedUntil,
		PasswordChangedAt: adminUser.PasswordChangedAt,
		TwoFactorEnabled:  adminUser.TwoFactorEnabled,
		TwoFactorSecret:   adminUser.TwoFactorSecret,
		CreatedAt:         adminUser.CreatedAt,
		UpdatedAt:         adminUser.UpdatedAt,
		Role:              role,
	}
}

func (s *AdminService) mapToAdminRoleResponse(adminRole *models.AdminRoles) *web.AdminRoleResponse {
	return &web.AdminRoleResponse{
		Id:          adminRole.Id,
		Code:        adminRole.Code,
		Name:        adminRole.Name,
		Description: adminRole.Description,
		IsActive:    adminRole.IsActive,
		CreatedAt:   adminRole.CreatedAt,
		UpdatedAt:   adminRole.UpdatedAt,
	}
}

func (s *AdminService) updateAdminRole(adminRole **models.AdminRoles, payload *web.AdminRoleResponse) {
	//if (*adminRole).Code == payload.Code {
	//	(*adminRole).Code = payload.Code
	//}

	if (*adminRole).Name != payload.Name {
		(*adminRole).Name = payload.Name
	}

	if (*adminRole).Description != payload.Description {
		(*adminRole).Description = payload.Description
	}

	if (*adminRole).IsActive != payload.IsActive {
		(*adminRole).IsActive = payload.IsActive
	}
}

func (s *AdminService) updateAdminUser(adminUser **models.AdminUsers, payload *web.DetailAdminResponse) {
	if (*adminUser).Email != payload.Email {
		(*adminUser).Email = payload.Email
	}

	if (*adminUser).FullName != payload.FullName {
		(*adminUser).FullName = payload.FullName
	}

	if (*adminUser).Phone != payload.Phone {
		(*adminUser).Phone = payload.Phone
	}

	if (*adminUser).RoleId != payload.RoleId {
		(*adminUser).RoleId = payload.RoleId
	}

	if (*adminUser).IsActive != payload.IsActive {
		(*adminUser).IsActive = payload.IsActive
	}

	if (*adminUser).IsVerified != payload.IsVerified {
		(*adminUser).IsVerified = payload.IsVerified
	}
}
