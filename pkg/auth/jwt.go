package auth

import (
	"share-docs/pkg/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

var (
	accessTokenSecret  = []byte(util.MustGetEnv("JWT_ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(util.MustGetEnv("JWT_REFRESH_TOKEN_SECRET"))
)

func GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error) {
	accessTokenClaims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "share-docs",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenSigned, err := accessToken.SignedString(accessTokenSecret)

	if err != nil {
		return nil, err
	}

	refreshTokenClaims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "refresh_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "share-docs",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)
	refreshTokenSigned, err := refreshToken.SignedString(refreshTokenSecret)

	if err != nil {
		return nil, err
	}

	tokenPair := &TokenPair{
		AccessToken:  accessTokenSigned,
		RefreshToken: refreshTokenSigned,
	}

	return tokenPair, nil
}
