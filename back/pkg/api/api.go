package api

import (
	"backend/back/pkg/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	r  *mux.Router
	db *repository.PGRepo
}

func New(router *mux.Router, db *repository.PGRepo) *api {
	return &api{r: router, db: db}
}

func (api *api) Handle() {
	api.r.HandleFunc("/api/lectures", api.lectures)
	api.r.HandleFunc("/api/lectures/{id}", api.GetLectureById).Methods("GET")
	api.r.HandleFunc("/api/search", api.SearchLecture).Methods("POST")
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
