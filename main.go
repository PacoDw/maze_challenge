package main

import (
	"os"

	"github.com/PacoDw/maze_challenge/repository"
	"github.com/PacoDw/maze_challenge/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	router.Use(
		repository.GinMiddleware(os.Getenv("MOGODB_CONN")),
	)

	spot := router.Group("/spot")
	{
		spot.POST("/create", routes.CreateSpot)
		spot.GET("/read/:id", routes.GetSpot)
		spot.PATCH("/update", routes.UpdateSpot)
		spot.DELETE("/delete/:id", routes.DeleteSpot)
	}

	quadrant := router.Group("/quadrant")
	{
		quadrant.POST("/create", routes.CreateQuadrant)
		quadrant.GET("/read/:id", routes.GetQuadrant)
		quadrant.PATCH("/update", routes.UpdateQuadrant)
		quadrant.DELETE("/delete/:id", routes.DeleteQuadrant)
	}

	router.Run(":3000")
}
