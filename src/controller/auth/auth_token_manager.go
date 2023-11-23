package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
)

type UnauthorizedError struct {
	message string
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		message: message,
	}
}

func (e *UnauthorizedError) Error() string {
	return e.message
}

type TokenSet struct {
	AccessToken  string
	RefreshToken string
}

type AppUserClaims struct {
	LoginID   string `json:"loginId"`
	AppUserID int    `json:"appUserId"`
	Username  string `json:"username"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

type authTokenManager struct {
	signingKey     []byte
	signingMethod  jwt.SigningMethod
	tokenTimeout   time.Duration
	refreshTimeout time.Duration
}

type AuthTokenManager interface {
	CreateTokenSet(ctx context.Context, appUser domain.AppUserModel) (*TokenSet, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
}

func NewAuthTokenManager(signingKey []byte, signingMethod jwt.SigningMethod, tokenTimeout, refreshTimeout time.Duration) AuthTokenManager {
	return &authTokenManager{
		signingKey:     signingKey,
		signingMethod:  signingMethod,
		tokenTimeout:   tokenTimeout,
		refreshTimeout: refreshTimeout,
	}
}

func (m *authTokenManager) CreateTokenSet(ctx context.Context, appUser domain.AppUserModel) (*TokenSet, error) {
	accessToken, err := m.createJWT(ctx, appUser, m.tokenTimeout, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.createJWT(ctx, appUser, m.refreshTimeout, "refresh")
	if err != nil {
		return nil, err
	}

	return &TokenSet{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *authTokenManager) createJWT(ctx context.Context, appUser domain.AppUserModel, duration time.Duration, tokenType string) (string, error) {
	logger := liblog.GetLoggerFromContext(ctx, log.AppGatewayLoggerContextKey)
	now := time.Now()
	claims := AppUserClaims{
		LoginID:   appUser.GetLoginID(),
		AppUserID: appUser.GetAppUserID().Int(),
		Username:  appUser.GetUsername(),
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	logger.DebugContext(ctx, fmt.Sprintf("claims: %+v", claims))

	token := jwt.NewWithClaims(m.signingMethod, claims)
	signed, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", liberrors.Errorf(". err: %w", err)
	}

	return signed, nil
}

func (m *authTokenManager) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	logger := liblog.GetLoggerFromContext(ctx, log.AppGatewayLoggerContextKey)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return m.signingKey, nil
	}

	currentToken, err := jwt.ParseWithClaims(tokenString, &AppUserClaims{}, keyFunc)
	if err != nil {
		logger.InfoContext(ctx, "", slog.Any("err", err))
		return "", NewUnauthorizedError(fmt.Sprintf("failed to ParseWithClaims. err: %v", err))
	}

	currentClaims, ok := currentToken.Claims.(*AppUserClaims)
	if !ok || !currentToken.Valid {
		return "", NewUnauthorizedError("Invalid token")
	}

	if currentClaims.TokenType != "refresh" {
		return "", NewUnauthorizedError("Invalid token type")
	}

	now := time.Now()
	tmpID := 1
	userModel, err := libdomain.NewBaseModel(1, now, now, tmpID, tmpID)
	if err != nil {
		return "", liberrors.Errorf("libdomain.NewBaseModel. err: %w", err)
	}

	appUserID, err := domain.NewAppUserID(currentClaims.AppUserID)
	if err != nil {
		return "", liberrors.Errorf("domain.NewAppUserID. err: %w", err)
	}

	appUser, err := domain.NewAppUserModel(userModel, appUserID, currentClaims.LoginID, currentClaims.Username)
	if err != nil {
		return "", liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	accessToken, err := m.createJWT(ctx, appUser, m.tokenTimeout, "access")
	if err != nil {
		return "", liberrors.Errorf("m.createJWT. err: %w", err)
	}

	return accessToken, nil
}
