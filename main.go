package main

import (
	"fmt"
	"os"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/handler"
	"github.com/ButterHost69/odoo-hackathon/utils"
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

	var email_id string
	if os.Getenv("GMAIL_ID") == "" {
		fmt.Print("Enter Your Email ID: ")
		fmt.Scan(&email_id)
	}
	var password string
	if os.Getenv("GMAIL_PASS") == "" {
		fmt.Print("Enter Your Email ID Password: ")
		fmt.Scan(&password)
	}

	utils.InitEmailClient(email_id, password)
	err = utils.SMTP_SendMessagetoEmail("parzival1520@gmail.com", "Dummy Message", "Works!!")
	if err != nil {
		fmt.Println("[Log] Error In Sending Email !!")
		fmt.Println(err)
	}

	fmt.Println("[Log] Init Done !!")
}

func main() {
	fmt.Println("Hello World")
	ipAddr := "localhost:3030"

	fmt.Println("Server Running on: ", ipAddr)
	r := gin.Default()

	r.GET("/", handler.RenderInitPage)

	r.POST("/register", handler.CreateCompany)
	r.POST("/login", handler.Login)

	r.POST("/approval/:managerEmail/:expenseID/:status", handler.ApproveExpense)

	r.Run(ipAddr)
}
