package service

import (
	"cc-service1/storage"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestResponse struct {
	Status    string `json:"status" xml:"status"`
	RequestID string `json:"request-id" xml:"request-id"`
}

type StatusResponse struct {
	Status   string `json:"name" xml:"name"`
	ImageURL string `json:"image-url" xml:"image-url"`
}

func (s Service1) RequestHandler(e echo.Context) error {
	emailAddress := e.FormValue("email")
	file, err := e.FormFile("image")
	if err != nil {
		fmt.Println("Error at image file receive %w", err)
		return err
	}
	imageSrc, err := file.Open()
	if err != nil {
		fmt.Println("Error at image file open %w", err)
		return err
	}
	defer imageSrc.Close()

	requestID, err := s.DataBase.GenerateUniqueID()
	if err != nil {
		fmt.Println("Error at id generation: %w", err)
		return err
	}

	dbEntry := storage.PicRequestEntry{
		ReqId:        requestID,
		Email:        emailAddress,
		ReqStatus:    "pending",
		ImageCaption: "",
		NewImageURL:  "",
	}

	err = s.DataBase.Save(dbEntry)
	if err != nil {
		fmt.Println("Error at database insert: %w", err)
		return err
	}

	_, err = s.PicStore.Upload(imageSrc, file.Size, requestID)
	if err != nil {
		fmt.Println("Error at image upload: %w", err)
		return err
	}

	// TODO add request to message broker

	fmt.Printf("request registered successfuly with id: %s\n", requestID)
	return e.JSON(http.StatusAccepted, &RequestResponse{Status: "accepted", RequestID: requestID})
}

func (s Service1) StatusHandler(e echo.Context) error {
	// TODO
	return nil
}
