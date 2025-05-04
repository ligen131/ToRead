package controllers

// import (
// 	"errors"
// 	"strconv"
// 	"to-read/controllers/auth"
// 	"to-read/model"
// 	"to-read/utils/logs"

// 	"github.com/labstack/echo/v4"
// 	"go.uber.org/zap"
// 	"gorm.io/gorm"
// )

// func FindLocation(c echo.Context, request model.Location) (location model.Location, err error, isInternalServerError bool) {
// 	location = request
// 	if request.ID != 0 {
// 		location, err = model.FindLocationByID(request.ID)
// 	} else if request.LocationName != "" {
// 		location, err = model.FindLocationByName(request.LocationName)
// 	} else {
// 		return location, errors.New("Location ID or name is required."), false
// 	}
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return location, err, false
// 		}
// 		return location, ResponseInternalServerError(c, "Find location failed", err), true
// 	}
// 	return location, nil, false
// }

// type LocationGETResponse struct {
// 	ID           uint32 `json:"location_id"   `
// 	LocationName string `json:"location_name" `
// }

// func LocationGET(c echo.Context) error {
// 	logs.Debug("GET /location")

// 	locationRequest := model.Location{}
// 	num, _ := strconv.ParseUint(c.QueryParam("location_id"), 10, 32)
// 	locationRequest.ID = uint32(num)
// 	locationRequest.LocationName = c.QueryParam("location_name")

// 	location, err, e500 := FindLocation(c, model.Location{
// 		ID:           locationRequest.ID,
// 		LocationName: locationRequest.LocationName,
// 	})
// 	if e500 {
// 		return err
// 	}
// 	if err != nil {
// 		return ResponseBadRequest(c, "Find location failed.", err)
// 	}

// 	return ResponseOK(c, LocationGETResponse{
// 		ID:           location.ID,
// 		LocationName: location.LocationName,
// 	})
// }

// type LocationScanPOSTRequest struct {
// 	LocationID   uint32 `json:"location_id"   `
// 	LocationName string `json:"location_name" `
// 	UserID       uint32 `json:"user_id"`
// 	OpenID       string `json:"openid"`
// }

// func LocationScanPOST(c echo.Context) error {
// 	logs.Debug("POST /location/scan")

// 	locationRequest := LocationScanPOSTRequest{}
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
// 		return ResponseForbidden(c, "You cannot help other to scan.", nil)
// 	}

// 	scanned, err := model.IsUserScannedLocation(user.ID, location.ID)
// 	if err != nil {
// 		return ResponseInternalServerError(c, "Failed to find is_user_scanned_location.", err)
// 	}

// 	if !scanned {
// 		err = model.AddUserScanLocation(user.ID, location.ID)
// 		if err != nil {
// 			return ResponseInternalServerError(c, "Failed to add user_scan_location.", err)
// 		}
// 	}

// 	return ResponseOK(c, StatusMessage{
// 		Status: "ok",
// 	})
// }

// type UserScannedLocationListResponse struct {
// 	ScannedLocationList []LocationGETResponse `json:"scanned_location_list"`
// }

// func LocationListGET(c echo.Context) error {
// 	logs.Debug("GET /location/list")

// 	userRequest := model.User{}
// 	num, _ := strconv.ParseUint(c.QueryParam("user_id"), 10, 32)
// 	userRequest.ID = uint32(num)
// 	userRequest.OpenID = c.QueryParam("openid")

// 	user, err, e500 := FindUser(c, model.User{
// 		ID:     userRequest.ID,
// 		OpenID: userRequest.OpenID,
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
// 		return ResponseForbidden(c, "You cannot get other's scanned location list.", nil)
// 	}

// 	list, err := model.GetUserScannedLocationList(user.ID)
// 	if err != nil {
// 		return ResponseInternalServerError(c, "Get user scanned location list failed.", err)
// 	}

// 	resp := UserScannedLocationListResponse{
// 		ScannedLocationList: make([]LocationGETResponse, 0),
// 	}
// 	for _, usl := range list {
// 		location, err := model.FindLocationByID(usl.LocationID)
// 		if err != nil {
// 			logs.Warn("Find location for user_scanned_location failed.", zap.Uint32("LocationID", usl.LocationID), zap.Error(err))
// 		}
// 		resp.ScannedLocationList = append(resp.ScannedLocationList, LocationGETResponse{
// 			ID:           location.ID,
// 			LocationName: location.LocationName,
// 		})
// 	}

// 	return ResponseOK(c, resp)
// }
