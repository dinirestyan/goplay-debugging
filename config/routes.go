package config

import (
	"github.com/dinirestyan/goplay-debugging/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	// controller package
	// import helper
)

// CORS is function for enable cors on backend
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Routes is struct to store connection
type Routes struct {
	DB *gorm.DB
}

// Setup is function to setup the application routes
func (r *Routes) Setup() *gin.Engine {
	app := gin.Default()

	app.Use(CORS())
	main := app.Group("/")
	{
		goplayCtrl := controllers.GoplayController{DB: r.DB}
		goplay := main.Group("goplay")
		{
			goplay.POST("/login", goplayCtrl.Login)
			goplay.POST("/upload", goplayCtrl.Upload)
			goplay.GET("/list", goplayCtrl.GetFileLists)
		}
	}

	app.Use(cors.Default())
	return app
}
