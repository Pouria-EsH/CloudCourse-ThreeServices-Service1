package main

import (
	"cc-service1/broker"
	"cc-service1/service"
	"cc-service1/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}

	acs3_bucket := os.Getenv("CCSERV1_ACS3_BUCKET")
	acs3_region := os.Getenv("CCSERV1_ACS3_REGION")
	acs3_endpoint := os.Getenv("CCSERV1_ACS3_ENDPOINT")
	acs3_accessKey := os.Getenv("CCSERV1_ACS3_ACCESSKEY")
	acs3_secretKey := os.Getenv("CCSERV1_ACS3_SECRETKEY")
	imagestore, err := storage.NewArvanCloudS3(
		acs3_bucket,
		acs3_region,
		acs3_endpoint,
		acs3_accessKey,
		acs3_secretKey)
	if err != nil {
		log.Fatalf("Fatal error at object storage: %v\n", err)
	}

	mySQL_username := os.Getenv("CCSERV1_MYSQL_USERNAME")
	mySQL_password := os.Getenv("CCSERV1_MYSQL_PASSWORD")
	mySQL_address := os.Getenv("CCSERV1_MYSQL_ADDRESS")
	database, err := storage.NewMySQLDB(mySQL_username, mySQL_password, mySQL_address, "ccp1")
	if err != nil {
		log.Fatalf("Fatal error at database: %v\n", err)
	}

	cloudamq_url := os.Getenv("CCSERV1_AMQP_URL")
	cloudamq := broker.NewCloudAMQ(cloudamq_url, "cc-pr")

	srv := service.NewService1(*database, *imagestore, *cloudamq)
	err = srv.Execute()

	if err != nil {
		log.Fatalf("could not start service1: %v\n", err)
	}
}
