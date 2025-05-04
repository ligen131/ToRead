package controllers

// import (
// 	"strconv"
// 	"time"
// 	"to-read/controllers/auth"
// 	"to-read/model"
// 	"to-read/utils/logs"

// 	"github.com/labstack/echo/v4"
// 	"go.uber.org/zap"
// 	"gorm.io/gorm"
// )

// type PostCreateRequest struct {
// 	AuthorID     uint32 `json:"user_id"`
// 	AuthorName   string `json:"user_name"`
// 	AuthorOpenID string `json:"openid"`
// 	LocationID   uint32 `json:"location_id"`
// 	LocationName string `json:"location_name"`
// 	Content      string `json:"content"`
// }

// type PostCreateResponse struct {
// 	Status   string `json:"status"`
// 	PostID   uint32 `json:"post_id"`
// 	IsPublic bool   `json:"is_public"`
// }

// func PostPOST(c echo.Context) error {
// 	logs.Debug("POST /post")

// 	postRequest := PostCreateRequest{}
// 	_ok, err := Bind(c, &postRequest)
// 	if !_ok {
// 		return err
// 	}

// 	user, err, e500 := FindUser(c, model.User{
// 		ID:     postRequest.AuthorID,
// 		OpenID: postRequest.AuthorOpenID,
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

// 	if user.UserName != postRequest.AuthorName {
// 		model.UpdateUserName(user.ID, postRequest.AuthorName)
// 	}

// 	claims, err := auth.GetClaimsFromHeader(c)
// 	if err != nil {
// 		return ResponseBadRequest(c, err.Error(), nil)
// 	}
// 	if claims.ID != user.ID || claims.OpenID != user.OpenID {
// 		return ResponseForbidden(c, "You cannot post other's post.", nil)
// 	}

// 	location, err, e500 := FindLocation(c, model.Location{
// 		ID:           postRequest.LocationID,
// 		LocationName: postRequest.LocationName,
// 	})
// 	if e500 {
// 		return err
// 	}
// 	if err != nil {
// 		return ResponseBadRequest(c, "Find location failed.", err)
// 	}

// 	scanned, err := model.IsUserScannedLocation(user.ID, location.ID)
// 	if err != nil {
// 		return ResponseInternalServerError(c, "Failed to find is_user_scanned_location.", err)
// 	}
// 	if !scanned {
// 		return ResponseForbidden(c, "抱歉，您尚未前往该地点扫码，留言失败。", nil)
// 	}

// 	post, err := model.CreatePost(user.ID, location.ID, postRequest.AuthorName, time.Now(), postRequest.Content, true)
// 	if err != nil {
// 		return ResponseInternalServerError(c, "Failed to create post into database.", err)
// 	}

// 	return ResponseOK(c, PostCreateResponse{
// 		Status:   "留言成功！",
// 		PostID:   post.ID,
// 		IsPublic: post.IsPublic,
// 	})
// }

// type PostGetRequest struct {
// 	AuthorID            uint32 `json:"user_id"`
// 	AuthorOpenID        string `json:"openid"`
// 	LocationID          uint32 `json:"location_id"`
// 	LocationName        string `json:"location_name"`
// 	Limit               int    `json:"limit"`
// 	OrderBy             string `json:"order_by"`
// 	StartTime           int64  `json:"start_time"`
// 	IsIncludeRecentPost bool   `json:"is_include_recent_post"` // Currently useless
// }

// type PostResponse struct {
// 	AuthorID     uint32 `json:"user_id"`
// 	AuthorName   string `json:"user_name"`
// 	AuthorOpenID string `json:"openid"`
// 	LocationID   uint32 `json:"location_id"`
// 	LocationName string `json:"location_name"`
// 	PostID       uint32 `json:"post_id"`
// 	Time         int64  `json:"time"`
// 	Content      string `json:"content"`
// 	IsPublic     bool   `json:"is_public"`
// }

// type PostGetResponse struct {
// 	PostList []PostResponse `json:"post_list"`
// }

// func PostGET(c echo.Context) error {
// 	logs.Debug("GET /post")

// 	postRequest := PostGetRequest{}
// 	num, _ := strconv.ParseUint(c.QueryParam("user_id"), 10, 32)
// 	postRequest.AuthorID = uint32(num)
// 	postRequest.AuthorOpenID = c.QueryParam("openid")
// 	num, _ = strconv.ParseUint(c.QueryParam("location_id"), 10, 32)
// 	postRequest.LocationID = uint32(num)
// 	postRequest.LocationName = c.QueryParam("location_name")
// 	numint, _ := strconv.ParseInt(c.QueryParam("limit"), 10, 32)
// 	postRequest.Limit = int(numint)
// 	postRequest.OrderBy = c.QueryParam("order_by")
// 	numint, _ = strconv.ParseInt(c.QueryParam("start_time"), 10, 64)
// 	postRequest.StartTime = numint
// 	postRequest.IsIncludeRecentPost, _ = strconv.ParseBool(c.QueryParam("is_include_recent_post"))

// 	user, err, e500 := FindUser(c, model.User{
// 		ID:     postRequest.AuthorID,
// 		OpenID: postRequest.AuthorOpenID,
// 	})
// 	if e500 {
// 		return err
// 	}
// 	if err == gorm.ErrRecordNotFound {
// 		return ResponseBadRequest(c, "User not found.", err)
// 	} else if err != nil {
// 		user.ID = 0
// 	}

// 	location, err, e500 := FindLocation(c, model.Location{
// 		ID:           postRequest.LocationID,
// 		LocationName: postRequest.LocationName,
// 	})
// 	if e500 {
// 		return err
// 	}
// 	if err == gorm.ErrRecordNotFound {
// 		return ResponseBadRequest(c, "Location not found.", err)
// 	} else if err != nil {
// 		location.ID = 0
// 	}

// 	userMp := make(map[uint32]model.User)
// 	if user.ID != 0 {
// 		userMp[user.ID] = user
// 	}

// 	locationMp := make(map[uint32]model.Location)
// 	if location.ID != 0 {
// 		locationMp[location.ID] = location
// 	}

// 	limitUsed := 0
// 	excludePostID := uint32(0)
// 	resp := PostGetResponse{
// 		PostList: make([]PostResponse, 0),
// 	}

// 	claims, err := auth.GetClaimsFromHeader(c)

// 	if claims.ID != 0 {
// 		posts, err := model.GetPostsList(claims.ID, location.ID, time.Now(), false, "time", 1, 0)
// 		if err != nil {
// 			return ResponseInternalServerError(c, "Get posts list failed.", err)
// 		}

// 		for _, post := range posts {
// 			if userMp[post.AuthorID].ID == 0 {
// 				userMp[post.AuthorID], err = model.FindUserByID(post.AuthorID)
// 				if err != nil {
// 					logs.Warn("Find user for post failed.", zap.Uint32("AuthorID", post.AuthorID), zap.Error(err))
// 				}
// 			}

// 			if locationMp[post.LocationID].ID == 0 {
// 				locationMp[post.LocationID], err = model.FindLocationByID(post.LocationID)
// 				if err != nil {
// 					logs.Warn("Find location for post failed.", zap.Uint32("LocationID", post.LocationID), zap.Error(err))
// 				}
// 			}

// 			limitUsed += 1
// 			excludePostID = post.ID

// 			resp.PostList = append(resp.PostList, PostResponse{
// 				AuthorID:     post.AuthorID,
// 				AuthorName:   post.AuthorName,
// 				AuthorOpenID: userMp[post.AuthorID].OpenID,
// 				LocationID:   post.LocationID,
// 				LocationName: locationMp[post.LocationID].LocationName,
// 				PostID:       post.ID,
// 				Time:         post.Time.Unix(),
// 				Content:      post.Content,
// 				IsPublic:     post.IsPublic,
// 			})
// 		}
// 	}

// 	postRequest.Limit -= limitUsed
// 	if postRequest.Limit < 0 {
// 		postRequest.Limit = 0
// 	}

// 	posts, err := model.GetPostsList(user.ID, location.ID, time.Unix(postRequest.StartTime, 0), false, postRequest.OrderBy, postRequest.Limit, excludePostID)
// 	if err != nil {
// 		return ResponseInternalServerError(c, "Get posts list failed.", err)
// 	}

// 	for _, post := range posts {
// 		if userMp[post.AuthorID].ID == 0 {
// 			userMp[post.AuthorID], err = model.FindUserByID(post.AuthorID)
// 			if err != nil {
// 				logs.Warn("Find user for post failed.", zap.Uint32("AuthorID", post.AuthorID), zap.Error(err))
// 			}
// 		}

// 		if locationMp[post.LocationID].ID == 0 {
// 			locationMp[post.LocationID], err = model.FindLocationByID(post.LocationID)
// 			if err != nil {
// 				logs.Warn("Find location for post failed.", zap.Uint32("LocationID", post.LocationID), zap.Error(err))
// 			}
// 		}

// 		resp.PostList = append(resp.PostList, PostResponse{
// 			AuthorID:     post.AuthorID,
// 			AuthorName:   post.AuthorName,
// 			AuthorOpenID: userMp[post.AuthorID].OpenID,
// 			LocationID:   post.LocationID,
// 			LocationName: locationMp[post.LocationID].LocationName,
// 			PostID:       post.ID,
// 			Time:         post.Time.Unix(),
// 			Content:      post.Content,
// 			IsPublic:     post.IsPublic,
// 		})
// 	}

// 	return ResponseOK(c, resp)
// }
