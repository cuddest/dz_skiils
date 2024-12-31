package main

import (
	"log"

	"github.com/cuddest/dz-skills/config"
	"github.com/cuddest/dz-skills/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not start the application: %v", err)
	}
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	routes.InitRoutes(router)
	log.Println("Server running on port 8080...")
	router.Run(":8080")
}
