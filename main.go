package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dinirestyan/goplay-debugging/controllers"
	"github.com/dinirestyan/goplay-debugging/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var port string
	flag.StringVar(&port, "port", os.Getenv("PORT"), "port of the service")
	fmt.Println("port : " + os.Getenv("PORT"))
	db := utils.GetDBConnection()
	defer db.Close()

	router := gin.Default()
	goplay := router.Group("goplay")
	goplayCtrl := controllers.GoplayController{DB: db}
	{
		goplay.POST("/login", goplayCtrl.Login)

	}

	router.Run()

	// config := &config.Routes{DB: db}

	// routes := config.Setup()
	// routes.Run(fmt.Sprintf(":%s", port))
}
