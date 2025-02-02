package notification

import (
	"encoding/json"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"log"
	"net/http"
	"noti-share-backend/internal/helper"
	"noti-share-backend/internal/server"
)

type PostNotificationBody struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	OwnerId     string `json:"owner_id"`
	Owner       string `json:"owner"`
	Origin      string `json:"origin"`
	Timestamp   uint64 `json:"timestamp"`
	Version     uint8  `json:"version"`
	Destination string `json:"destination"`
	DestinationPlatform string `json:"destination_platform"`
}

func PostNotificationHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := PostNotificationBody{}
	err := decoder.Decode(&body)

	if err != nil {
		log.Println("Could not parse JSON: ", err)

		response := struct {
			Status string `json:"status"`
			Reason string `json:"reason"`
		}{Status: "Failure", Reason: "MALFORMED_JSON"}

		helper.ResponseJson(writer, 400, response)

		return
	}

	data := map[string]string{
		"title":     body.Title,
		"content":   body.Content,
		"owner_id":  body.OwnerId,
		"owner":     body.Owner,
		"origin":    body.Origin,
		"timestamp": fmt.Sprintf("%d", body.Timestamp),
		"version":   fmt.Sprintf("%d", body.Version),
	}

	var message *messaging.Message

	if body.DestinationPlatform == "iOS" {
		log.Println("Sending message to iOS...")

		message = &messaging.Message{
			APNS: buildAPNSConfig(body),
			Notification: buildAPNSNotification(body),
			Data: data,
			Token: body.Destination,
		}
	} else {
		log.Println("Sending message to Android...")

		message = &messaging.Message{
			Data: data,
			Token: body.Destination,
		}
	}
	
	_, err = resources.FirebaseMessaging.Send(request.Context(), message)

	if err != nil {
		log.Println("Error sending message:", err)

		errString := fmt.Sprintf("%v", err)

		if errString == "Requested entity was not found." || errString == "The registration token is not a valid FCM registration token" || errString == "exactly one of token, topic or condition must be specified" {
			response := struct {
				Status       string `json:"status"`
				Reason       string `json:"reason"`
				InvalidToken string `json:"invalid_token"`
			}{Status: "Failure", Reason: "INVALID_TOKEN", InvalidToken: body.Destination}

			helper.ResponseJson(writer, 404, response)
		} else {
			response := struct {
				Status        string `json:"status"`
				Reason        string `json:"reason"`
				FirebaseError string `json:"firebase_error"`
			}{Status: "Failure", Reason: "FIREBASE_ERROR", FirebaseError: errString}

			helper.ResponseJson(writer, 500, response)
		}
	} else {
		helper.ResponseNoContent(writer)
	}
}

func buildAPNSConfig(body PostNotificationBody) *messaging.APNSConfig{
	return &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: body.Owner,
					Body: body.Content,
					SubTitle: body.Title,
				},
			},
		},
	}
}

func buildAPNSNotification(body PostNotificationBody) *messaging.Notification {
	return &messaging.Notification{
		Title: body.Title,
		Body: body.Content,
	}
}
