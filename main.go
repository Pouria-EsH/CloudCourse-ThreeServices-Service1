package main

import (
	"cc-service1/service"
	"cc-service1/storage"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	acs3_accessKey := os.Getenv("CCSERV1_ACS3_ACCESSKEY")
	acs3_secretKey := os.Getenv("CCSERV1_ACS3_SECRETKEY")
	imagestore, err := storage.NewArvanCloudS3("cc-practice-004", "ir-thr-at1",
		"https://s3.ir-thr-at1.arvanstorage.com",
		acs3_accessKey,
		acs3_secretKey)
	if err != nil {
		fmt.Println("Fatal error at object storage: %w", err)
		os.Exit(1)
	}

	mySQL_username := os.Getenv("CCSERV1_MYSQL_USERNAME")
	mySQL_password := os.Getenv("CCSERV1_MYSQL_PASSWORD")
	database, err := storage.NewMySQLDB(mySQL_username, mySQL_password, "127.0.0.1:3306", "ccp1")
	if err != nil {
		fmt.Println("Fatal error at database: %w", err)
		os.Exit(1)
	}

	srv := service.NewService1(*database, *imagestore)
	err = srv.Execute()

	if err != nil {
		fmt.Println("Could not start service1")
	}
}
