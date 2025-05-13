package controllers

import (
	"net/http"
	"to-read/utils/logs"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ErrorMessage struct {
	Message string `json:"msg"`
	Err     string `json:"err"`
}

type StatusMessage struct {
	Status string `json:"status"`
}

type ResponseStruct struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Bind(c echo.Context, obj interface{}) (bool, error) {
	err := c.Bind(&obj)
	if err != nil {
		logs.Warn("Failed to parse request data.", zap.Error(err))
		return false, c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Failed to parse request data.",
				Err:     err.Error(),
			},
		})
	}
	logs.Debug("Parsed struct:", zap.Any("obj", obj))
	return true, nil
}

func BindJSON(c echo.Context, obj interface{}) (bool, error) {
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	err := c.Bind(&obj)

	if err != nil {
		logs.Warn("Failed to parse request data.", zap.Error(err))
		return false, c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Failed to parse request data.",
				Err:     err.Error(),
			},
		})
	}
	logs.Debug("Parsed struct:", zap.Any("obj", obj))
	return true, nil
}

func ResponseOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func ResponseBadRequest(c echo.Context, errMessage string, err error) error {
	Err := ""
	if err != nil {
		Err = err.Error()
	}
	return c.JSON(http.StatusBadRequest, ResponseStruct{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
		Data: ErrorMessage{
			Message: errMessage,
			Err:     Err,
		},
	})
}

func ResponseInternalServerError(c echo.Context, errMessage string, err error) error {
	Err := ""
	if err != nil {
		Err = err.Error()
	}
	return c.JSON(http.StatusInternalServerError, ResponseStruct{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Data: ErrorMessage{
			Message: errMessage,
			Err:     Err,
		},
	})
}

func ResponseUnauthorized(c echo.Context, errMessage string, err error) error {
	Err := ""
	if err != nil {
		Err = err.Error()
	}
	return c.JSON(http.StatusUnauthorized, ResponseStruct{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
		Data: ErrorMessage{
			Message: errMessage,
			Err:     Err,
		},
	})
}

func ResponseForbidden(c echo.Context, errMessage string, err error) error {
	Err := ""
	if err != nil {
		Err = err.Error()
	}
	return c.JSON(http.StatusForbidden, ResponseStruct{
		Code:    http.StatusForbidden,
		Message: "Forbidden",
		Data: ErrorMessage{
			Message: errMessage,
			Err:     Err,
		},
	})
}
