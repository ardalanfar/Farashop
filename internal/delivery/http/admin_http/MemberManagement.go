package admin_http

import (
	"Farashop/internal/adapter/store"
	"Farashop/internal/contract"
	"Farashop/internal/dto"
	"Farashop/internal/service/admin_service"
	"Farashop/pkg/customerror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ShowMembers(conn store.DbConn) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err      error
			request  dto.ShowMembersRequest
			response dto.ShowMembersResponse
		)

		//send service
		response, err = admin_service.NewAdmin(conn).ShowMembers(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, customerror.Unsuccessful())
		}
		//return ui
		return ctx.JSON(http.StatusOK, response)
	}
}

func DeleteMember(conn store.DbConn, validator contract.ValidateDeleteMember) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err      error
			request  dto.DeleteMemberRequest
			response dto.DeleteMemberResponse
		)

		//bind user information
		userID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		request = dto.DeleteMemberRequest{ID: uint(userID)}
		//validat information
		err = validator(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		//service
		response, err = admin_service.NewAdmin(conn).DeleteMember(ctx.Request().Context(), request)
		if err != nil && !response.Result {
			return echo.NewHTTPError(http.StatusInternalServerError, customerror.InfoIncorrect())
		}

		//return ui
		return ctx.JSON(http.StatusOK, nil)
	}
}

func ShowInfoMember(conn store.DbConn, validator contract.ValidateShowInfoMember) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err      error
			request  dto.ShowInfoMemberRequest
			response dto.ShowInfoMemberResponse
		)

		//bind user information
		userID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		request = dto.ShowInfoMemberRequest{ID: uint(userID)}

		//validat information
		err = validator(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		//service
		response, err = admin_service.NewAdmin(conn).ShowInfoMember(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		//return ui
		return ctx.JSON(http.StatusOK, response)
	}
}
