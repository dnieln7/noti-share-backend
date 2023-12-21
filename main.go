package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"noti-share-backend/internal/env"
	"noti-share-backend/internal/handler/notification"
	"noti-share-backend/internal/server"
	"time"
)

func main() {
	properties := env.GetEnvProperties()
	resources := buildResources()

	router := buildRouter(resources)
	httpServer := buildHttpServer(properties, router)

	log.Println("Server will start on port: ", properties.Port)

	err := httpServer.ListenAndServe()

	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
}

func buildResources() *server.Resources {
	app, err := firebase.NewApp(context.Background(), nil)

	if err != nil {
		log.Fatal("Could not build firebase app: ", err)
	}

	messaging, err := app.Messaging(context.Background())

	if err != nil {
		log.Fatal("Could not build messaging client: ", err)
	}

	return &server.Resources{
		Firebase:          app,
		FirebaseMessaging: messaging,
	}
}

func buildRouter(resources *server.Resources) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/notifications", resources.HttpHandler(notification.PostNotificationHandler)).
		Methods("POST")

	return router
}

func buildHttpServer(properties *env.EvnProperties, router *mux.Router) *http.Server {
	httpServer := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + properties.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return httpServer
}
