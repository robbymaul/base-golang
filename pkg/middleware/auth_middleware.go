package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/config"
	pkgjwt "paymentserviceklink/pkg/jwt"
	"strings"
)

type Auth struct {
	cfg     *config.Config
	adapter *pkgjwt.JwtAdapter
	repo    *repositories.RepositoryContext
}

func NewAuth(cfg *config.Config, repo *repositories.RepositoryContext) *Auth {
	adapter := pkgjwt.NewJwtAdapter(cfg.JwtIssuer, cfg.JwtSecret)
	return &Auth{
		cfg:     cfg,
		adapter: adapter,
		repo:    repo,
	}
}

func (a *Auth) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authType, v := a.getAuthorizationHeaderValue(ctx, header)
		log.Info().Msg(fmt.Sprintf("[authentication] auth type = [%v] with value = [%v]", authType, v))
		if authType == "" && v == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		switch authType {
		case Basic:
			ctxAuth := a.getValueBasicAuth(ctx, v)
			log.Debug().Msg(fmt.Sprintf("[case basic] [get value basic auth] contect auth = [%v]", ctxAuth))

			basicAuth := a.getBasicAuth(ctxAuth)
			if basicAuth == nil {
				log.Warn().Msg(fmt.Sprintf("[case basic] [get value basic auth] basic auth = [%v]", basicAuth))
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			match := a.matchUsernamePassword(basicAuth)
			if !match {
				log.Warn().Msg(fmt.Sprintf("[case basic] [match username password value basic auth] match = [%v] with context = [%v] ",
					match,
					basicAuth),
				)
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.Next()
		case Bearer:
			ctxToken := a.getValueBearerToken(ctx, v)

			token := a.getBearerToken(ctxToken)

			if token == "" {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			validation := a.tokenValidation(token)
			if validation == nil {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			ctx.Set(BearerToken, validation)
			ctx.Next()
		default:
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func (a *Auth) getAuthorizationHeaderValue(ctx *gin.Context, header string) (authType string, value string) {
	getHeader := ctx.GetHeader(header)
	log.Info().Msg(fmt.Sprintf("[get authorization header value] with value = [%v]", getHeader))
	if getHeader == "" {
		return "", ""
	}

	temp := strings.SplitN(getHeader, " ", 2)
	if len(temp) != 2 {
		return "", ""
	}

	return temp[0], temp[1]
}

func (a *Auth) getValueBearerToken(ctx context.Context, token string) context.Context {
	log.Info().Str("token", token).Msg(fmt.Sprintf("[get value bearer token] with value = [%v]", token))

	if token == "" {
		return ctx
	}

	return context.WithValue(ctx, Bearer, token)
}

func (a *Auth) getBearerToken(ctx context.Context) string {
	s := ctx.Value(Bearer).(string)
	log.Info().Str("s", s).Msg(fmt.Sprintf("[get value bearer token] with value = [%v]", s))

	return s
}

func (a *Auth) tokenValidation(token string) *pkgjwt.JwtResponse {
	jwt, err := a.adapter.VerifyJwt(token)
	if err != nil {
		log.Error().Err(err).Msg("verify jwt token error")
		return nil
	}

	return jwt
}
