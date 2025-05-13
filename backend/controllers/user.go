package controllers

import (
	"regexp"
	"strconv"
	"to-read/controllers/auth"
	"to-read/model"
	"to-read/utils"
	"to-read/utils/logs"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserLoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	ID                  uint32 `json:"user_id"`
	UserName            string `json:"user_name"`
	AccessToken         string `json:"token"`
	AccessTokenExpireAt int64  `json:"token_expiration_time"`
}

func UserLoginPOST(c echo.Context) error {
	logs.Debug("POST /user/login")

	userRequest := UserLoginRequest{}
	_ok, err := Bind(c, &userRequest)
	if !_ok {
		return err
	}

	user, err := model.FindUserByName(userRequest.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResponseBadRequest(c, "User or password incorrect.", nil)
		}
		return ResponseInternalServerError(c, "Find user failed.", err)
	}
	if user.Deleted {
		return ResponseBadRequest(c, "User or password incorrect.", nil)
	}
	if user.PasswordMD5 != utils.GetMD5(userRequest.Password) {
		return ResponseBadRequest(c, "User or password incorrect.", nil)
	}

	accessTokenString, accessTokenExpireAt, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return ResponseInternalServerError(c, "Generate access token failed.", err)
	}

	return ResponseOK(c, UserLoginResponse{
		ID:                  user.ID,
		UserName:            user.UserName,
		AccessToken:         accessTokenString,
		AccessTokenExpireAt: accessTokenExpireAt.Unix(),
	})
}

type UserRegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	ID       uint32 `json:"user_id"   `
	UserName string `json:"user_name" `
	Role     uint32 `json:"role"      `
}

func isAlphanumeric(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
	return match
}

func UserRegisterPOST(c echo.Context) error {
	logs.Debug("POST /user/register")

	userRequest := UserRegisterRequest{}
	_ok, err := Bind(c, &userRequest)
	if !_ok {
		return err
	}

	if !isAlphanumeric(userRequest.UserName)	 {
		return ResponseBadRequest(c, "User name must be alphanumeric.", nil)
	}

	user, err := model.UserRegister(userRequest.UserName, userRequest.Password)
	if err != nil {
		return ResponseInternalServerError(c, "Register user failed.", err)
	}

	return ResponseOK(c, UserRegisterResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Role:     user.Role,
	})
}

func UserIsAuthGET(c echo.Context) error {
	logs.Debug("GET /user/isauth")

	return ResponseOK(c, StatusMessage{
		Status: "OK",
	})
}

type UserGETResponse struct {
	ID       uint32 `json:"user_id"   `
	UserName string `json:"user_name" `
	Role     uint32 `json:"role"      `
}

func UserGET(c echo.Context) error {
	logs.Debug("GET /user")

	var err error
	user := model.User{}
	num, _ := strconv.ParseUint(c.QueryParam("user_id"), 10, 32)
	userID := uint32(num)
	userName := c.QueryParam("user_name")

	if userID != 0 {
		user, err = model.FindUserByID(userID)
	} else if userName != "" {
		user, err = model.FindUserByName(userName)
	} else {
		return ResponseBadRequest(c, "User ID or user_name is required.", nil)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResponseBadRequest(c, "User not found.", nil)
		}
		return ResponseInternalServerError(c, "Find user failed", err)
	}

	return ResponseOK(c, UserGETResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Role:     user.Role,
	})
}
