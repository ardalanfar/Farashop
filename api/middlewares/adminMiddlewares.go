package middlewares

import (
	"Farashop/internal/entity"
	"Farashop/pkg/auth"
	"Farashop/pkg/customerror"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetAdminGroup(grp *echo.Group) {

	grp.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &auth.Claims{},
		SigningKey:              []byte(auth.GetJWTSecret()),
		TokenLookup:             "cookie:access-token", // "<source>:<name>"
		ErrorHandlerWithContext: auth.JWTErrorChecker,
	}))

	grp.Use(TokenRefresherMiddlewareAdmin)
	grp.Use(CheckAccessAdmin)
}

func CheckAccessAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Get("user") == nil {
			return next(ctx)
		}
		u := ctx.Get("user").(*jwt.Token)
		claims := u.Claims.(*auth.Claims)
		if int(claims.Access) == 1 {
			return next(ctx)
		}
		return ctx.JSON(http.StatusBadRequest, customerror.NOAccess())
	}
}

func TokenRefresherMiddlewareAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Get("user") == nil {
			return next(ctx)
		}
		u := ctx.Get("user").(*jwt.Token)
		claims := u.Claims.(*auth.Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < (10 * time.Minute) {
			rc, err := ctx.Cookie(auth.GetrefReshTokenCookieName())

			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(auth.GetRefreshJWTSecret()), nil
				})

				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						ctx.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					_ = auth.GenerateTokensAndSetCookies(entity.User{
						Username: claims.Name,
						ID:       claims.ID,
						Access:   claims.Access,
					}, ctx)
				}
			}
		}
		return next(ctx)
	}
}
