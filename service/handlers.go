package service

import (
	"cc-service1/storage"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestResponse struct {
	Status    string `json:"status" xml:"status"`
	RequestID string `json:"request-id" xml:"request-id"`
}

type StatusResponse struct {
	Status   string `json:"status" xml:"name"`
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
		fmt.Println("Error at database insert:", err)
		return err
	}

	_, err = s.PicStore.Upload(imageSrc, file.Size, requestID)
	if err != nil {
		fmt.Println("Error at image upload:", err)
		s.failureHandler(requestID)
		return err
	}

	err = s.MsgBroker.Send(requestID)
	if err != nil {
		fmt.Println("Error at rabbitMQ send:", err)
		s.failureHandler(requestID)
		return err
	}

	fmt.Printf("request registered successfuly with id: %s\n", requestID)
	return e.JSON(http.StatusAccepted, &RequestResponse{Status: "accepted", RequestID: requestID})
}

func (s Service1) StatusHandler(e echo.Context) error {
	requestId := e.FormValue("request-id")
	reqEntry, err := s.DataBase.Get(requestId)
	if err != nil {
		var notfound *storage.RequestNotFoundError
		if errors.As(err, &notfound) {
			return e.JSON(http.StatusNotFound, &StatusResponse{
				Status:   "Not Found",
				ImageURL: "",
			})
		}
		s.failureHandler(requestId)
		return err
	}

	fmt.Printf("request for status with id: %s was found successfuly\n", requestId)

	return e.JSON(http.StatusOK, &StatusResponse{
		Status:   reqEntry.ReqStatus,
		ImageURL: reqEntry.NewImageURL,
	})
}

func (s Service1) failureHandler(requstId string) {
	err := s.DataBase.SetStatus(requstId, "failure")
	if err != nil {
		var notfound *storage.RequestNotFoundError
		if !errors.As(err, &notfound) {
			fmt.Printf("couldn't update request %s status to \"failed\": %v\n", requstId, err)
		}
		return
	}
	fmt.Printf("request %s status is set to 'failure'", requstId)
}
