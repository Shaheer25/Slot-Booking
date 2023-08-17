package routes

import (
	controller "github.com/Shaheer25/go-auth/controllers"
	"github.com/Shaheer25/go-auth/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/users", controller.GetUsers())

	incomingRoutes.GET("/users/:user_id", controller.GetUser())

	incomingRoutes.POST("/users/availabletime", controller.AvailabilityTime())

	incomingRoutes.POST("/tickets/generate", controller.GenerateTickets())

	incomingRoutes.DELETE("/tickets/delete", controller.DeleteTicket())

	incomingRoutes.GET("/admin/reservations", controller.GetAllReservations())

}
