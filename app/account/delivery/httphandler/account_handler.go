package httphandler

import (
	"login-system/server/domains"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

type RespFormat struct {
	Message string `json:"message"`
}

type AccountHandler struct {
	AccountUsecase domains.AccountUsecase
}

func NewAccountHandler(e *echo.Echo, ua domains.AccountUsecase) {
	handler := &AccountHandler{
		AccountUsecase: ua,
	}

	e.POST("/api/v1/signup", handler.SignUp)
}

func (h *AccountHandler) SignUp(c echo.Context) error {
	var account domains.Account
	err := c.Bind(&account)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, RespFormat{Message: err.Error()})
	}

	var ok bool
	if ok, err = isRequestValid(&account); !ok {
		return c.JSON(http.StatusBadRequest, RespFormat{Message: err.Error()})
	}

	ctx := c.Request().Context()

	// check existed email
	_, err = h.AccountUsecase.GetByEmail(ctx, account.Email)
	if err != nil && err != domains.ErrNotFound {
		return c.JSON(getStatusCode(err), RespFormat{Message: err.Error()})
	}

	// store account to db
	err = h.AccountUsecase.Store(ctx, &account)
	if err != nil {
		return c.JSON(getStatusCode(err), RespFormat{Message: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

func isRequestValid(m *domains.Account) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domains.ErrInternalServerError:
		return http.StatusInternalServerError
	case domains.ErrNotFound:
		return http.StatusNotFound
	case domains.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
