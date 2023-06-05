package member_http

import (
	"Farashop/internal/adapter/store"
	"Farashop/internal/dto"
	"Farashop/internal/service/member_service"
	"Farashop/pkg/auth"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func ShowOrders(conn store.DbConn) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			response dto.ShowOrdersResponse
		)

		u := ctx.Get("user").(*jwt.Token)
		claims := u.Claims.(*auth.Claims)

		req := dto.ShowOrdersRequest{
			ID: claims.ID,
		}

		//Service
		response, err := member_service.NewMember(conn).ShowOrders(ctx.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		//return ui
		return ctx.JSON(http.StatusOK, response)
	}
}
