package main

import (
	"fmt"
	"os"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env/main.env")

	err := db.InitDB()
	if err != nil {
		fmt.Println("[main.Init] Failed To Init Postgress Connection : \n", err.Error())
		os.Exit(1)
	}

	fmt.Println("[Log] Init Done !!")
}

func main() {
	fmt.Println("Hello World")
	ipAddr := "localhost:3030"

	fmt.Println("Server Running on: ", ipAddr)
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		handler.RenderAuthPage(ctx, "")
	})

	r.POST("/register", handler.CreateCompany)
	r.POST("/login", handler.Login)

	r.Run(ipAddr)
}
