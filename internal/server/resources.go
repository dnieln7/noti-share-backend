package server

import (
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"net/http"
)

type Resources struct {
	Firebase          *firebase.App
	FirebaseMessaging *messaging.Client
}

type ResourcesHandlerFunc func(writer http.ResponseWriter, request *http.Request, resources *Resources)

func (resources *Resources) HttpHandler(handlerFunc ResourcesHandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handlerFunc(writer, request, resources)
	}
}
