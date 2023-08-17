package routes

import (
	controller "github.com/Shaheer25/go-auth/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("users/signup", controller.Signup())

	incomingRoutes.POST("users/login", controller.Login())

	incomingRoutes.POST("/generate-refresh-token", controller.GenerateRefreshToken())

	incomingRoutes.POST("users/logout", controller.Logout())
}
