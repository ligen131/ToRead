package controllers

// import (
// 	"time"
// 	"to-read/controllers/auth"
// 	"to-read/model"
// 	"to-read/utils/logs"

// 	"github.com/labstack/echo/v4"
// )

// type SharePOSTRequest struct {
// 	UserID       uint32 `json:"user_id"`
// 	OpenID       string `json:"openid"`
// 	LocationID   uint32 `json:"location_id"   `
// 	LocationName string `json:"location_name" `
// }

// func SharePOST(c echo.Context) error {
// 	logs.Debug("POST /share")

// 	locationRequest := SharePOSTRequest{}
// 	_ok, err := Bind(c, &locationRequest)
// 	if !_ok {
// 		return err
// 	}

// 	location, err, e500 := FindLocation(c, model.Location{
// 		ID:           locationRequest.LocationID,
// 		LocationName: locationRequest.LocationName,
// 	})
// 	if e500 {
// 		return err
// 	}
// 	if err != nil {
// 		return ResponseBadRequest(c, "Find location failed.", err)
// 	}

// 	user, err, e500 := FindUser(c, model.User{
// 		ID:     locationRequest.UserID,
// 		OpenID: locationRequest.OpenID,
// 	})
// 	if e500 {
// 		return err
// 	}
// 	if err != nil {
// 		return ResponseBadRequest(c, "Find user failed.", err)
// 	}

// 	if user.Deleted {
// 		return ResponseBadRequest(c, "This user has been deleted.", nil)
// 	}

// 	claims, err := auth.GetClaimsFromHeader(c)
// 	if err != nil {
// 		return ResponseBadRequest(c, err.Error(), nil)
// 	}
// 	if claims.ID != user.ID || claims.OpenID != user.OpenID {
// 		return ResponseForbidden(c, "You cannot help other to share.", nil)
// 	}

// 	err = model.AddShareRecord(user.ID, location.ID, time.Now())
// 	if err != nil {
// 		return ResponseInternalServerError(c, "Failed to add share record.", err)
// 	}

// 	return ResponseOK(c, StatusMessage{
// 		Status: "ok",
// 	})
// }
