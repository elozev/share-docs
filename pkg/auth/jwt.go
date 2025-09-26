package auth

import (
	"errors"
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

const (
	accessTokenExpiration  = 24 * time.Hour
	refreshTokenExpiration = 7 * 24 * time.Hour
)

func RefreshAccessToken(c Claims) (*string, error) {
	accessTokenClaims := &Claims{
		UserID:    c.UserID,
		Email:     c.Email,
		TokenType: "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "share-docs",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenSigned, err := accessToken.SignedString(accessTokenSecret)

	if err != nil {
		return nil, err
	}

	return &accessTokenSigned, err
}

func GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error) {
	accessTokenClaims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "share-docs",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiration)),
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

type Token int

const (
	AccessToken Token = iota
	RefreshToken
)

func ValidateToken(tokenString string, tokenType Token) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if tokenType == AccessToken {
			return accessTokenSecret, nil
		} else if tokenType == RefreshToken {
			return refreshTokenSecret, nil
		} else {
			return nil, errors.New("Unsupported JWT token type")
		}
	})

	if err != nil {
		return nil, err
	}

	return claims, err
}
