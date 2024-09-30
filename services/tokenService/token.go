package tokenservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	refreshtoken "github.com/Jesuloba-world/koodle-server/services/tokenService/refreshToken"
)

type TokenService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	refreshtokens   RefreshTokenRepository
}

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(refreshToken, userId string, expiresAt time.Time) error
	RevokeRefreshToken(refreshToken string) error
	IsRefreshTokenRevoked(refreshToken string) (bool, error)
}

func NewTokenService(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration, db *bun.DB) *TokenService {
	refreshTokenRepo := refreshtoken.NewRefreshTokenRepo(db)

	return &TokenService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		refreshtokens:   refreshTokenRepo,
	}
}

func (t *TokenService) GenerateTokens(userID string) (string, string, error) {
	accessToken, err := t.generateToken(userID, AccessTokenType, t.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.generateToken(userID, RefreshTokenType, t.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	// store refresh token
	expiresAt := time.Now().Add(t.refreshTokenTTL)
	err = t.refreshtokens.StoreRefreshToken(refreshToken, userID, expiresAt)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t *TokenService) generateToken(userID, tokenType string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    userID,
		"exp":    time.Now().Add(ttl).Unix(),
		"type":   tokenType,
		"random": uuid.New().String(),
	})

	signedToken, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return signedToken, nil
}

func (t *TokenService) ValidateToken(token string) (jwt.MapClaims, string, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return t.secretKey, nil
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse token: %w", err)
	}

	tokentype, ok := claims["type"].(string)
	if !ok {
		return nil, "", errors.New("invalid token type")
	}

	// check if refresh token is revoked
	if tokentype == RefreshTokenType {
		refreshToken := token
		revoked, err := t.refreshtokens.IsRefreshTokenRevoked(refreshToken)
		if err != nil {
			return nil, "", err
		}
		if revoked {
			return nil, "", errors.New("refresh token has been revoked")
		}
	}

	return claims, tokentype, nil
}

func (t *TokenService) RevokeRefreshToken(refreshToken string) error {
	return t.refreshtokens.RevokeRefreshToken(refreshToken)
}

func (t *TokenService) RefreshToken(oldRefreshToken string) (string, string, error) {
	// validate the old refresh token
	claims, tokenType, err := t.ValidateToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	if tokenType != RefreshTokenType {
		return "", "", errors.New("invalid token type")
	}

	// revoke the old refresh token
	err = t.RevokeRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	// extract user id from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID in token claims")
	}

	// generate new access and refresh tokens
	return t.GenerateTokens(userID)
}
