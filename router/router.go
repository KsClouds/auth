package router

import (
	"auth/api"

	"github.com/gin-gonic/gin"
)

func StartRouter() (router *gin.Engine) {
	router = gin.Default()
	router.Use(Cors())

	router.GET("/", api.Welcome)
	user := router.Group("user")
	user.GET("/checkStatus", api.CheckStatus)
	user.POST("/login", api.Login)
	user.POST("/register", api.Register)
	user.DELETE("/logout", api.Logout)

	user.Use(AuthRequired())
	user.PUT("password", api.ChangePassword)
	user.PUT("repassword/:id", api.ResetPassword)
	user.PUT("superior/:id/:superiorId", api.ChangeSuperior)
	user.DELETE("manage/:id", api.DeleteUser)
	user.GET("list", api.GetList)
	user.GET("list/all", api.GetListAll)
	user.GET("me", api.UserMe)
	user.GET("subList", api.GetSubList)
	return
}
