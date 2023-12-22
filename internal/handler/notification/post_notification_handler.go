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
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Owner        string   `json:"owner"`
	Origin       string   `json:"origin"`
	Timestamp    uint64   `json:"timestamp"`
	Version      uint8    `json:"version"`
	Destinations []string `json:"destinations"`
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

	var invalidTokens []string

	for _, token := range body.Destinations {
		message := &messaging.Message{
			Data: map[string]string{
				"title":     body.Title,
				"content":   body.Content,
				"owner":     body.Owner,
				"origin":    body.Origin,
				"timestamp": fmt.Sprintf("%d", body.Timestamp),
				"version":   fmt.Sprintf("%d", body.Version),
			},
			Token: token,
		}

		_, err := resources.FirebaseMessaging.Send(request.Context(), message)

		if err != nil {
			log.Println("Error sending message:", err)

			errString := fmt.Sprintf("%v", err)

			if errString == "Requested entity was not found." || errString == "The registration token is not a valid FCM registration token" || errString == "exactly one of token, topic or condition must be specified" {
				invalidTokens = append(invalidTokens, token)
			}
		} else {
			log.Println("Success sending message")
		}
	}

	if len(invalidTokens) > 0 {
		response := struct {
			Status        string   `json:"status"`
			Reason        string   `json:"reason"`
			InvalidTokens []string `json:"invalid_tokens"`
		}{Status: "Failure", Reason: "INVALID_TOKEN", InvalidTokens: invalidTokens}

		helper.ResponseJson(writer, 400, response)
	} else {
		helper.ResponseNoContent(writer)
	}
}
