package constant

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shaheer25/go-auth/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CodeRunner() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Err loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "4856"

	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.AdminRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted for api-1",
		})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted for api-2",
		})
	})

	defer TrackTime(time.Now())

	time.Sleep(time.Second)
	router.Run(":" + port)

}
func TrackTime(t time.Time) {
	elapsed := time.Since(t)
	fmt.Printf("Execution completed in %v\n", elapsed)
}
