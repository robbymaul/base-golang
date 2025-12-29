package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

type BasicAuthRequest struct {
	Username string
	Password string
}

func (a *Auth) getBasicAuth(ctx context.Context) *BasicAuthRequest {
	value, ok := ctx.Value(BasicAuth).(BasicAuthRequest)
	log.Debug().Msg(fmt.Sprintf("[get basic auth] value = [%v] and ok = [%v]", value, ok))
	if !ok {
		return nil
	}

	return &value
}

func (a *Auth) matchUsernamePassword(basicAuth *BasicAuthRequest) bool {
	if basicAuth.Username == "" || basicAuth.Password == "" {
		return false
	} else if basicAuth.Username != a.cfg.BasicAuthUsername || basicAuth.Password != a.cfg.BasicAuthPassword {
		return false
	} else {
		return true
	}
}

func (a *Auth) getValueBasicAuth(ctx context.Context, v string) context.Context {
	log.Debug().Msg(fmt.Sprintf("[get value basic auth] value v string = [%s]", v))

	decodeString, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("[decode string base 64] [get value basic auth] error = [%v]", err))
		return ctx
	}

	valueString := string(decodeString)
	log.Debug().Msg(fmt.Sprintf("[get value basic auth] value string (decode string) = [%s]", valueString))

	indexByte := strings.IndexByte(valueString, ':')
	log.Debug().Msg(fmt.Sprintf("[get value basic auth] index byte (value string) = [%v]", valueString))
	if indexByte < 0 {
		return ctx
	}

	username, password := valueString[:indexByte], valueString[indexByte+1:]

	return context.WithValue(ctx, BasicAuth, BasicAuthRequest{username, password})
}
