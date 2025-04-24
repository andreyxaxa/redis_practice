package server

import (
	"encoding/json"
	"net/http"
	"os"
	"red_test/internal/app/models"
	"red_test/internal/app/storage"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router  *mux.Router
	storage storage.Storage
}

func New(storage storage.Storage) *Server {
	return &Server{
		router:  mux.NewRouter(),
		storage: storage,
	}
}

func (s *Server) Start() error {
	s.router.HandleFunc("/users/create", s.handleCreateUser).Methods("POST")
	s.router.HandleFunc("/users/get/{id}", s.handleGetUser).Methods("GET")

	httpPort := ":" + os.Getenv("HTTP_PORT")
	return http.ListenAndServe(httpPort, s.router)
}

// handlers ...

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.User().CreateUser(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	u, err := s.storage.User().GetUser(id)
	if err != nil {
		if err == redis.Nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
