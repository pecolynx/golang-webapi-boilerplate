package middleware

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/controller/auth"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
)

func NewAuthMiddleware(signingKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := liblog.GetLoggerFromContext(ctx, log.AppTraceLoggerContextKey)

		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			logger.WarnContext(ctx, "invalid header. Bearer not found")
			return
		}

		tokenString := authorization[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &auth.AppUserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil {
			logger.WarnContext(ctx, "invalid token", slog.Any("err", err))
			return
		}

		if claims, ok := token.Claims.(*auth.AppUserClaims); ok && token.Valid {
			c.Set("AuthorizedUser", claims.AppUserID)

			logger.InfoContext(ctx, "", slog.String("uri", c.Request.RequestURI), slog.Int("operator_id", claims.AppUserID))
		} else {
			logger.WarnContext(ctx, "invalid token")
			return
		}
	}
}
