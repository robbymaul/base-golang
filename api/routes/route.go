package routes

import (
	"application/app/controllers"
	"application/app/repositories"
	"application/config"
	"application/pkg/middleware"
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

	// r.router.GET("", ctrl.HealthController)
	r.initRoute()
}

func (r *Route) initRoute() {

}
