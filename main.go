package main

import (
	"fmt"
	"os"

	"github.com/ButterHost69/odoo-hackathon/db"
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

func main(){
	fmt.Println("Hello World")
}