package routes

import (
	"paymentserviceklink/app/controllers"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Route struct {
	startTime  time.Time
	appVersion string
	cfg        *config.Config
	repo       *repositories.RepositoryContext
	router     *gin.RouterGroup
	auth       *middleware.Auth
	ctrl       *controllers.Controller
}

func NewRoute(startTime time.Time, appVersion string, cfg *config.Config, repo *repositories.RepositoryContext, router *gin.RouterGroup) *Route {
	return &Route{startTime: startTime, appVersion: appVersion, cfg: cfg, repo: repo, router: router}
}

func (r *Route) RegisterCoreServicesRoutes() {
	ctrl := controllers.NewController(r.startTime, r.appVersion, r.cfg, r.repo)
	log.Warn().Msg("running route ....")

	auth := middleware.NewAuth(r.cfg, r.repo)

	if r.auth == nil {
		r.auth = auth
	}

	if r.ctrl == nil {
		r.ctrl = ctrl
	}

	r.router.GET("", ctrl.HealthController)
	r.initRoute()
}

func (r *Route) initRoute() {
	r.authRoute()
	r.adminRoute()
	r.aggregatorRoute()
	r.configurationRoute()
	r.platformRoute()
	r.paymentMethodRoute()
	r.paymentRoute()
	r.webhook()
	r.espayRoute()
	r.midtransRoute()
	r.transactionRoute()
	r.kWalletRoute()
	r.assetRoute()
}

func (r *Route) authRoute() {
	r.router.POST("/auth/login", r.auth.Authentication(), r.ctrl.AdminLoginController)
	r.router.GET("/auth/me", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminMeController)
}

func (r *Route) adminRoute() {
	r.adminUserRoute()
	r.adminRoleRoute()
}

func (r *Route) adminRoleRoute() {
	r.router.POST("/admins/roles", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.CreateAdminRoleController)
	r.router.GET("/admins/roles", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.GetListAdminRoleController)
	r.router.GET("/admins/roles/:id", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.GetDetailAdminRoleController)
	r.router.PUT("/admins/roles", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.UpdateAdminRoleController)
}

func (r *Route) adminUserRoute() {
	r.router.POST("/admins", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.CreateAdminUserController)
	r.router.GET("/admins", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.GetListAdminUserController)
	r.router.GET("/admins/:id", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.GetDetailAdminUserController)
	r.router.PUT("/admins/:id", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.UpdateAdminUserController)
	r.router.DELETE("/admins/:id", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.DeleteAdminUserController)
}

func (r *Route) platformRoute() {
	r.router.POST("/admins/platforms", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreatePlatformController)
	r.router.GET("/admins/platforms", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListPlatformController)
	r.router.GET("/admins/platforms/:id-platform", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetDetailPlatformController)
	r.router.PUT("/admins/platforms/:id-platform", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdatePlatformController)
	r.router.PUT("/admins/platforms/:id-platform/secret-key", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdatePlatformSecretKeyController)
	//r.router.POST("/admins/platforms/:id-platform/platform-configurations", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreatePlatformConfigurationController)
	r.router.GET("/admins/platforms/:id-platform/platform-configurations", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListConfigurationController)
	r.router.GET("/admins/platforms/:id-platform/platform-configurations/:id-platform-configuration", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetConfigurationController)
	r.router.POST("/admins/platforms/:id-platform/payment-methods", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreateChannelController)
	r.router.GET("/admins/platforms/:id-platform/payment-methods", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.GetListPlatformChannelController)
	r.router.GET("/admins/platforms/:id-platform/payment-methods/:id-payment-method", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetDetailChannelController)
	r.router.PUT("/admins/platforms/:id-platform/payment-methods", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdateChannelController)
	r.router.POST("/platforms/payment-methods", r.auth.ClientMiddleware(), r.ctrl.ClientGetListChannelController)
	r.router.POST("/admins/platforms/:id-platform/configurations/assignment", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminAssignmentPlatformConfigurationController)
	r.router.POST("/admins/platforms/:id-platform/configurations/removal", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminRemovalPlatformConfigurationController)
	r.router.GET("/admins/platforms/:id-platform/configurations", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListConfigurationInPlatformController)
	r.router.POST("/admins/platforms/:id-platform/channels/assignment", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminAssignmentPlatformChannelController)
	r.router.GET("/admins/platforms/:id-platform/channels/assignment", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminAssignmentPlatformListChannelController)
	r.router.POST("/admins/platforms/:id-platform/channels/removal", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminRemovalPlatformChannelController)
	r.router.GET("/admins/platforms/:id-platform/channels", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListChannelInPlatformController)
}

func (r *Route) paymentMethodRoute() {
	r.router.POST("/admins/channels", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreateChannelController)
	r.router.GET("/admins/channels", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListChannelController)
	r.router.GET("/admins/channels/:id-channel", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetDetailChannelController)
	r.router.PUT("/admins/channels/:id-channel", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdateChannelController)
	r.router.POST("/payment-methods", r.auth.ClientMiddleware(), r.ctrl.ClientGetListChannelController)
}

func (r *Route) paymentRoute() {
	r.router.POST("/payments", r.auth.ClientMiddleware(), r.ctrl.CreatePaymentController)
	r.router.POST("/payments/detail", r.auth.ClientMiddleware(), r.ctrl.GetDetailPaymentController)
	r.router.POST("/payments/status", r.auth.ClientMiddleware(), r.ctrl.CheckStatusPaymentController)
	r.router.GET("/check/midtrans", r.ctrl.CheckKeyMidtransController)
}

func (r *Route) webhook() {
	r.router.POST("/senangpay/callback-notification", r.ctrl.WebhookCallbackSenangpayPaymentNotification)
}

func (r *Route) aggregatorRoute() {
	r.router.POST("/admins/aggregators", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreateAggregatorController)
	r.router.GET("/admins/aggregators", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetAllAggregatorController)
	r.router.GET("/admins/aggregators/:id-aggregator", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetAggregatorController)
	r.router.PUT("/admins/aggregators/:id-aggregator", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdateAggregatorController)
	//r.router.POST("/admins/aggregators/:id-aggregator/channels", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreateChannelController)
	//r.router.GET("/admins/aggregators/:id-aggregator/channels", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListChannelController)
	//r.router.GET("/admins/aggregators/:id-aggregator/channels/:id-channel", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetDetailChannelController)
	//r.router.PUT("/admins/aggregators/:id-aggregator/channels/:id-channel", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdateChannelController)
}

func (r *Route) configurationRoute() {
	r.router.POST("/admins/configurations", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminCreateConfigurationController)
	r.router.GET("/admins/configurations", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListConfigurationController)
	r.router.GET("/admins/configurations/:id-configuration", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetConfigurationController)
	r.router.PUT("/admins/configurations/:id-configuration", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminUpdateConfigurationController)
	r.router.POST("/admins/configurations/:id-configuration/platforms/assignment", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminAssignmentConfigurationToPlatformController)
	r.router.POST("/admins/configurations/:id-configuration/platforms/removal", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminRemovalPlatformFromConfigurationController)
	r.router.GET("/admins/configurations/:id-configuration/platforms", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListPlatformInConfigurationController)
}

func (r *Route) espayRoute() {
	r.router.GET("/espay/payment/inquiry", r.ctrl.EspayValidationInquiryController)
	r.router.POST("/espay/payment/inquiry", r.ctrl.EspayValidationInquiryController)
	r.router.GET("/espay/payment/notification", r.ctrl.EspayPaymentNotificationController)
	r.router.POST("/espay/payment/notification", r.ctrl.EspayPaymentNotificationController)
	r.router.POST("/espay/topup/notification", r.ctrl.EspayTopupNotificationController)
}

func (r *Route) midtransRoute() {
	r.router.POST("/midtrans/payment/notification", r.ctrl.MidtransPaymentNotification)
}

func (r *Route) transactionRoute() {
	r.router.GET("/admins/transactions", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListTransactionController)
	r.router.GET("/admins/transactions/:id-transaction", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetDetailTransactionController)
}

func (r *Route) kWalletRoute() {
	r.router.GET("/admins/k-wallets/members", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListKWalletController)
	r.router.GET("/admins/k-wallets/members/:no-rekening", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetDetailKWalletController)
	r.router.GET("/admins/k-wallets/members/:no-rekening/transactions", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListKWalletTransactionController)
	r.router.GET("/admins/topups/k-wallets", r.auth.Authentication(), r.auth.AdminAuthorization(), r.ctrl.AdminGetListTopupKWalletController)
	r.router.POST("/k-wallets/members/registration", r.auth.ClientMiddleware(), r.ctrl.CreateKWalletController)
	r.router.POST("/k-wallets/members", r.auth.ClientMiddleware(), r.ctrl.GetKWalletMemberController)
	r.router.POST("/k-wallets/members/virtual-account", r.auth.ClientMiddleware(), r.ctrl.GetVirtualAccountKWalletController)
	r.router.POST("/k-wallets/members/transactions", r.auth.ClientMiddleware(), r.ctrl.GetListKWalletTransactionMemberController)
	r.router.POST("/k-wallets/members/create-topup", r.auth.ClientMiddleware(), r.ctrl.CreateTopupKWalletController)
}

func (r *Route) assetRoute() {
	r.router.POST("/assets/upload-file/image", r.auth.Authentication(), r.ctrl.AdminUploadFileImageController)
}
