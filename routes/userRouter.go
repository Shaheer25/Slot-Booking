package routes

import (
	controller "github.com/Shaheer25/go-auth/controllers"
	"github.com/Shaheer25/go-auth/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/users/showtickets", controller.ShowTickets())

	incomingRoutes.POST("/tickets/book/:id", controller.BookTickets())

	incomingRoutes.GET("/user/reservations", controller.GetUserReservations())

	incomingRoutes.DELETE("/reservations/delete/:id",controller.DeleteReservation())
	
	incomingRoutes.GET("/tickets/empty",controller.GetEmptySlots())
}
