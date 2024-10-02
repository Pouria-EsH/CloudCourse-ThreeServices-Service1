package service

import (
	"cc-service1/broker"
	"cc-service1/storage"

	"github.com/labstack/echo/v4"
)

type Service1 struct {
	DataBase  storage.RequestDB
	PicStore  storage.ImageStorage
	MsgBroker broker.CloudAMQ
}

func NewService1(db storage.RequestDB, imgstore storage.ImageStorage, msgbroker broker.CloudAMQ) *Service1 {
	return &Service1{
		DataBase:  db,
		PicStore:  imgstore,
		MsgBroker: msgbroker,
	}
}

func (s Service1) Execute() error {
	app := echo.New()
	app.POST("/request", s.RequestHandler)
	app.GET("/status", s.StatusHandler)
	return app.Start(":8080")
}
