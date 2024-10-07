package service

import (
	"cc-service1/storage"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request-id"`
}

type StatusResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image-url"`
}

func (s Service1) RequestHandler(e echo.Context) error {
	emailAddress := e.FormValue("email")
	file, err := e.FormFile("image")
	if err != nil {
		return fmt.Errorf("at image file receive: %w", err)
	}
	imageSrc, err := file.Open()
	if err != nil {
		return fmt.Errorf("at image file open: %w", err)
	}
	defer imageSrc.Close()

	requestID, err := s.DataBase.GenerateUniqueID()
	if err != nil {
		return fmt.Errorf("at id generation: %w", err)
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
		return fmt.Errorf("at database insert: %w", err)
	}

	_, err = s.PicStore.Upload(imageSrc, file.Size, requestID)
	if err != nil {
		s.failureHandler(requestID)
		return fmt.Errorf("at image upload: %w", err)
	}

	err = s.MsgBroker.Send(requestID)
	if err != nil {
		s.failureHandler(requestID)
		return fmt.Errorf("at rabbitMQ send: %w", err)
	}

	log.Println("request registered successfuly with id: ", requestID)
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

	log.Printf("request for status with id: %s was found successfuly\n", requestId)

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
			log.Printf("couldn't update request %s status to \"failed\": %v\n", requstId, err)
		}
		return
	}
	log.Printf("request %s status is set to 'failure'", requstId)
}
