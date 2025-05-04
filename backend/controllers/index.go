package controllers

import (
	"to-read/utils/logs"

	"github.com/labstack/echo/v4"
)

type documentLink struct {
	Doc string `json:"document"`
}

type link struct {
	Link documentLink `json:"link"`
}

func IndexGET(c echo.Context) error {
	logs.Debug("GET /")

	return ResponseOK(c, link{
		Link: documentLink{
			Doc: "https://github.com/ligen131/ToRead/blob/main/backend/docs/api.md",
		},
	})
}
