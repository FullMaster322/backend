package main

import (
	"backend/back/pkg/api"
	"backend/back/pkg/repository"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db, err := repository.New("postgres://postgres:psql@localhost:5432/msfbd")
	if err != nil {
		log.Fatal(err.Error())
	}

	router := mux.NewRouter()
	api := api.New(router, db)
	api.Handle()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
		}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Content-Length"}),
	)

	log.Fatal(http.ListenAndServe("localhost:1212", corsHandler(router)))
}
