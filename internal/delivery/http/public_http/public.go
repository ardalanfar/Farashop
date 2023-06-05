package public_http

import (
	"Farashop/internal/adapter/sendmsg"
	"Farashop/internal/adapter/store"
	"Farashop/internal/config"
	"Farashop/internal/contract"
	"Farashop/internal/dto"
	"Farashop/internal/service/public_service"
	"Farashop/pkg/auth"
	"Farashop/pkg/customerror"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(conn store.DbConn, validator contract.ValidateRegister) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err      error
			request  dto.RegisterUserRequest
			response dto.RegisterUserResponse
		)

		//bind information
		err = json.NewDecoder(ctx.Request().Body).Decode(&request)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		//validat information
		err = validator(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		//service
		response, err = public_service.NewPublic(conn).Register(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		//send email success register
		if response.Result == true {
			configSys := config.GetConfig().Email
			to := []string{
				request.Email,
			}
			subject := "verify register"
			request := sendmsg.Mail{
				Sender:  configSys.From,
				To:      to,
				Subject: subject,
			}
			//build message
			msg := request.BuildMessage()
			//send email
			err = request.SendEmailUser(msg, configSys)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		//return ui
		return ctx.JSON(http.StatusOK, nil)
	}
}

func Login(conn store.DbConn, validator contract.ValidateLogin) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err      error
			request  dto.LoginUserRequest
			response dto.LoginUserResponse
		)

		//bind information
		err = json.NewDecoder(ctx.Request().Body).Decode(&request)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		//validat information
		err = validator(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		//service
		response, err = public_service.NewPublic(conn).Login(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		//call pkg auth(create token and cookie)
		if response.Result == true {
			err := auth.GenerateTokensAndSetCookies(response.User, ctx)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, customerror.Unsuccessful())
			}
		}

		//return ui
		return echo.NewHTTPError(http.StatusOK, "Wellcom "+response.User.Username)
	}
}

func MemberValidation(conn store.DbConn, validator contract.ValidateMemberValidation) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err      error
			request  dto.MemberValidationRequest
			response dto.MemberValidationResponse
		)

		//bind user information
		err = json.NewDecoder(ctx.Request().Body).Decode(&request)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		//validat information
		err = validator(ctx.Request().Context(), request)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		//service
		response, err = public_service.NewPublic(conn).MemberValidation(ctx.Request().Context(), request)
		if err != nil || !response.Result == true {
			return echo.NewHTTPError(http.StatusInternalServerError, customerror.InfoIncorrect())
		}

		//return ui
		return ctx.JSON(http.StatusOK, nil)
	}
}
