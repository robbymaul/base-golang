package pkgjwt

import (
	"application/app/web"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// JwtAdapter handles JWT creation and management.
type JwtAdapter struct {
	Algorithm *jwt.SigningMethodHMAC
	Issuer    string
	Secret    string
}

// IssueJwtPayload represents the payload used to create a JWT.
type IssueJwtPayload struct {
	Id      int64
	Subject string
	// Role     enums.CodeAdminRole
	Lifetime int64
}

// JwtResponse represents the JWT response format.
type JwtResponse struct {
	Id   int64  `json:"id"`
	Sub  string `json:"sub"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
}

// NewJwtAdapter creates a new instance of JwtAdapter.
func NewJwtAdapter(issuer, secret string) *JwtAdapter {
	return &JwtAdapter{
		Algorithm: jwt.SigningMethodHS512,
		Issuer:    issuer,
		Secret:    secret,
	}
}

// IssueJwt issues a new JWT based on the provided payload.
func (j *JwtAdapter) IssueJwt(payload *IssueJwtPayload) (*web.Session, error) {
	exp := time.Now().Add(time.Hour * time.Duration(payload.Lifetime)).Unix()

	token := jwt.New(j.Algorithm)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = payload.Id
	claims["exp"] = exp
	claims["sub"] = payload.Subject
	claims["iss"] = j.Issuer // Adding issuer to claims for more context

	tokenString, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to issue JWT token: %w", err)
	}

	return &web.Session{Token: tokenString, ExpiredAt: exp}, nil
}

func (j *JwtAdapter) VerifyJwt(token string) (*JwtResponse, error) {

	key := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Err(errors.New("unexpected signing method"))
			return nil, errors.New(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]).Error())
		}

		return []byte(j.Secret), nil
	}

	mapClaims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, &mapClaims, key); err != nil {
		log.Error().Err(err).Msg("failed parse with claims to verify JWT token")
		return nil, Error(err)
	}

	jsonData, err := json.Marshal(mapClaims)
	if err != nil {
		log.Error().Err(err).Msg("failed json marshal to verify JWT token")
		return nil, Error(err)
	}

	var claims JwtResponse
	if errMarshal := json.Unmarshal(jsonData, &claims); errMarshal != nil {
		log.Error().Err(errMarshal).Msg("failed json unmarshal to verify JWT token")
		return nil, Error(errMarshal)
	}

	return &claims, nil
}
