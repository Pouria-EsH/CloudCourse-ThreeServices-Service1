package main

import (
	"cc-service1/broker"
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
	mySQL_address := "mysql-container:3306"
	database, err := storage.NewMySQLDB(mySQL_username, mySQL_password, mySQL_address, "ccp1")
	if err != nil {
		fmt.Println("Fatal error at database: %w", err)
		os.Exit(1)
	}

	cloudamq_url := os.Getenv("CCSERV1_AMQP_URL")
	cloudamq := broker.NewCloudAMQ(cloudamq_url, "cc-pr")

	srv := service.NewService1(*database, *imagestore, *cloudamq)
	err = srv.Execute()

	if err != nil {
		fmt.Println("Could not start service1")
	}
}
