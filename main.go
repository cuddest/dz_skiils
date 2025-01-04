package main

import (
	"log"

	"github.com/cuddest/dz-skills/config"
	_ "github.com/cuddest/dz-skills/docs"
	"github.com/cuddest/dz-skills/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title          DZ Skills API
// @version        1.0
// @description    API for managing teachers, students, courses, and more in the DZ Skills Online Teaching platform.
// @contact.name   Support Team
// @contact.email  a_touati@estin.dz

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host           https://dzskiils-production.up.railway.app/
// @BasePath
// @comment        GitHub Repository:Available soon

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not start the application: %v", err)
	}
	/* because if my dumb ass main crashes, IT IS ME WHO HAVE TO CLEAN THE CONNECTION TRASH LEFT HERE */
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to extract *sql.DB: %v", err)
			return
		}
		sqlDB.Close()
	}()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not extract *sql.DB from *gorm.DB: %v", err)
	}

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	routes.InitRoutes(router, sqlDB)

	log.Println("Server running on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
