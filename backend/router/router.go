package router

import (
	"to-read/controllers"
	"to-read/controllers/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Load(e *echo.Echo) {
	routes(e)
}

func routes(e *echo.Echo) {
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	apiVersionUrl := "/api/v1"

	e.GET(apiVersionUrl+"", controllers.IndexGET)
	e.GET(apiVersionUrl+"/", controllers.IndexGET)

	e.GET(apiVersionUrl+"/health", controllers.HealthGET)

	userGroup := e.Group(apiVersionUrl + "/user")
	{
		userGroup.GET("", controllers.UserGET)
		userGroup.GET("/", controllers.UserGET)
		userGroup.POST("/login", controllers.UserLoginPOST)
		userGroup.POST("/register", controllers.UserRegisterPOST)
		userGroup.GET("/isauth", controllers.UserIsAuthGET, middleware.TokenVerificationMiddleware)
	}

	collectionGroup := e.Group(apiVersionUrl + "/collection")
	{
		collectionGroup.GET("", controllers.CollectionListGET, middleware.TokenVerificationMiddleware)
		collectionGroup.GET("/", controllers.CollectionListGET, middleware.TokenVerificationMiddleware)
		collectionGroup.GET("/list", controllers.CollectionListGET, middleware.TokenVerificationMiddleware)
		collectionGroup.POST("/add", controllers.CollectionAddPOST, middleware.TokenVerificationMiddleware)
		collectionGroup.GET("/summary", controllers.CollectionSummaryGET, middleware.TokenVerificationMiddleware)
		collectionGroup.GET("/tag", controllers.CollectionTagGET, middleware.TokenVerificationMiddleware)
	}

	// postGroup := e.Group(apiVersionUrl + "/post")
	// {
	// 	postGroup.POST("", controllers.PostPOST, middleware.TokenVerificationMiddleware)
	// 	postGroup.POST("/", controllers.PostPOST, middleware.TokenVerificationMiddleware)
	// 	postGroup.GET("", controllers.PostGET)
	// 	postGroup.GET("/", controllers.PostGET)
	// }

	// locationGroup := e.Group(apiVersionUrl + "/location")
	// {
	// 	locationGroup.GET("", controllers.LocationGET)
	// 	locationGroup.GET("/", controllers.LocationGET)
	// 	locationGroup.POST("/scan", controllers.LocationScanPOST, middleware.TokenVerificationMiddleware)
	// 	locationGroup.GET("/list", controllers.LocationListGET, middleware.TokenVerificationMiddleware)
	// }

	// shareGroup := e.Group(apiVersionUrl + "/share")
	// {
	// 	shareGroup.POST("", controllers.SharePOST, middleware.TokenVerificationMiddleware)
	// 	shareGroup.POST("/", controllers.SharePOST, middleware.TokenVerificationMiddleware)
	// }
}
