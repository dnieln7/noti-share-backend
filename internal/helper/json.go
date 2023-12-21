package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseNoContent(writer http.ResponseWriter) {
	writer.WriteHeader(204)
	_, err := writer.Write([]byte{})

	if err != nil {
		log.Println("Error writing: ", err)
	}
}
func ResponseJson(writer http.ResponseWriter, code int, payload interface{}) {
	bytes, err := json.Marshal(payload)

	if err != nil {
		log.Println("Failed to marshal payload")

		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)

	_, err = writer.Write(bytes)

	if err != nil {
		log.Println("Error writing: ", err)
	}
}
