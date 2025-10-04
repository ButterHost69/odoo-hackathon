package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // registers driver
)

var db *sql.DB

func InitDB() error {
	fmt.Println("[Log] Connecting to Postgress")
	sqlPassword := os.Getenv("POSTGRES_PASSWORD")
	if sqlPassword == "" {
		fmt.Print("Enter Your My POSTGRES Database Password: ")
		fmt.Scan(&sqlPassword)
	}

	sqlUsername := os.Getenv("POSTGRES_USERNAME")
	sqlDBName := os.Getenv("POSTGRES_DBNAME")
	
	sqlDBIP := os.Getenv("POSTGRES_IP")
	sqlDBPort := os.Getenv("POSTGRES_PORT")

	dblink := fmt.Sprintf(
    "postgres://%s:%s@%s:%s/%s?sslmode=disable",
    	sqlUsername, sqlPassword, sqlDBIP, sqlDBPort, sqlDBName,
	)

	var err error
	db, err = sql.Open("postgres", dblink)
	if err != nil {
		fmt.Println("[db.InitDB] error in connecting to postgress db: ", err)
		return err
	}

	fmt.Println("[Log] Connected to postgress db is succesful !!")

	// Pinging The Database To Verify The Connection
	if err = db.Ping(); err != nil {
		fmt.Println("[db.InitDB] error in pinging to postgress db: ", err)
		return err
	}

	fmt.Println("[Log] Pinging to postgress db is succesful !!")

	return nil
}