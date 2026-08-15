package router

import (
	"encoding/json"
	"fmt"
	"stock-tracker/app/auth"
	"stock-tracker/app/base"
	"stock-tracker/app/dashboard"
	"stock-tracker/app/item"
	"stock-tracker/app/itemlog"
	"stock-tracker/app/user"
	"stock-tracker/config"
	"stock-tracker/constant"
	"stock-tracker/db"
	"stock-tracker/jwt"
	"stock-tracker/model"
	"stock-tracker/response"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"net/http"
	_ "stock-tracker/docs"
)

func Init() *echo.Echo {
	router := websiteRouter()

	// Repository
	userRepository := user.NewRepository()
	itemRepository := item.NewRepository()
	itemlogRepository := itemlog.NewRepository()

	// Usecase
	baseUsecase := base.NewUsecase()
	authUsecase := auth.NewUsecase(baseUsecase, userRepository)
	userUsecase := user.NewUsecase(baseUsecase, userRepository)
	itemUsecase := item.NewUsecase(baseUsecase, itemRepository)
	itemlogUsecase := itemlog.NewUsecase(baseUsecase, itemlogRepository, itemRepository)
	dashboardUsecase := dashboard.NewUsecase(baseUsecase, itemRepository)

	// Handler
	authHandler := auth.NewHandler(authUsecase)
	userHandler := user.NewHandler(userUsecase)
	itemHandler := item.NewHandler(itemUsecase)
	itemlogHandler := itemlog.NewHandler(itemlogUsecase)
	dashboardHandler := dashboard.NewHandler(dashboardUsecase)

	if config.Debug {
		router.GET("/", func(c echo.Context) error {
			return response.Success(http.StatusOK, "Welcome", nil).SendJSON(c)
		})
		router.GET("/swg/*", echoSwagger.WrapHandler)
	}

	routerAuth := router.Group("/auth")
	routerAuth.POST("/sign-in", authHandler.SignIn)
	routerAuth.POST("/sign-out", authHandler.SignOut)
	routerAuth.GET("/init", authHandler.Init, checkTokenMiddleware)
	routerAuth.GET("/refresh-token", authHandler.RefreshToken, checkTokenMiddleware)

	routerDashboard := router.Group("/dashboard", checkTokenMiddleware)
	routerDashboard.GET("", dashboardHandler.Dashboard)

	routerUser := router.Group("/user", checkTokenMiddleware)
	routerUser.GET("", userHandler.Page)
	routerUser.POST("", userHandler.Create)
	routerUser.POST("/change-password", userHandler.ChangePassword)
	routerUser.PUT("/:id", userHandler.Update)
	routerUser.GET("/:id", userHandler.GetById)
	routerUser.DELETE("/:id", userHandler.Delete)

	routerItem := router.Group("/item", checkTokenMiddleware)
	routerItem.GET("", itemHandler.Page)
	routerItem.POST("", itemHandler.Create)
	routerItem.PUT("/:id", itemHandler.Update)
	routerItem.GET("/:id", itemHandler.GetById)
	routerItem.DELETE("/:id", itemHandler.Delete)

	routerItemlog := router.Group("/itemlog", checkTokenMiddleware)
	routerItemlog.GET("", itemlogHandler.Page)
	routerItemlog.POST("", itemlogHandler.Create)
	routerItemlog.PUT("/:id", itemlogHandler.Update)
	routerItemlog.GET("/:id", itemlogHandler.GetById)
	routerItemlog.DELETE("/:id", itemlogHandler.Delete)

	return router

}

func httpErrorHandler(err error, c echo.Context) {
	var errorResponse *response.Response
	code := http.StatusInternalServerError
	switch e := err.(type) {
	case *echo.HTTPError:
		// Handle pada saat URL yang di request tidak ada. atau ada kesalahan server.
		code = e.Code
		errorResponse = &response.Response{
			Status:  false,
			Message: fmt.Sprintf("%v", e.Message),
			Code:    code,
		}
	case *response.Response:
		errorResponse = e
	default:
		// Handle error dari panic
		code = http.StatusInternalServerError
		if config.Debug {
			errorResponse = &response.Response{
				Status:  false,
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		} else {
			errorResponse = &response.Response{
				Status:  false,
				Message: response.ErrorInternalServer,
				Code:    http.StatusInternalServerError,
			}
		}
	}

	js, err := json.Marshal(errorResponse)
	if err == nil {
		_ = c.Blob(code, echo.MIMEApplicationJSON, js)
	} else {
		b := []byte("{status: false, code: 500, message: \"unresolved error\"}")
		_ = c.Blob(code, echo.MIMEApplicationJSON, b)
	}
}

func checkTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		userLogin, err := jwt.ExtractClaims(c.Request().Header.Get(constant.AuthHeaderKey))
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, err.Error()).SendJSON(c)
		}

		conn := db.GetPostgresConnection()

		var user model.User
		err = conn.Where("id = ? ", userLogin.UserID).First(&user).Error
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewareUserNotFound).SendJSON(c)
		}

		if user.PassVersion != userLogin.PassVersion {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewarePassVersion).SendJSON(c)
		}

		c.Set(constant.TokenUserContext, userLogin)
		return next(c)
	}
}
