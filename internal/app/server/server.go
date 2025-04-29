package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"red_test/internal/app/data/cache"
	rediscache "red_test/internal/app/data/cache/redis_cache"
	"red_test/internal/app/data/storage"
	postgrestorage "red_test/internal/app/data/storage/postgre_storage"
	"red_test/internal/app/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type server struct {
	router  *mux.Router
	storage storage.Storage
	cache   cache.Cache
}

func newServer(storage storage.Storage, cache cache.Cache) *server {
	s := &server{
		router:  mux.NewRouter(),
		storage: storage,
		cache:   cache,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users/create", s.handleCreateUser).Methods("POST")
	s.router.HandleFunc("/users/get/{id}", s.handleGetUser).Methods("GET")
}

func Start(databaseURL string, cacheAddr string) error {
	db, err := newDB(databaseURL)
	if err != nil {
		return err
	}
	storage := postgrestorage.New(db)

	redis := newRedis(cacheAddr)
	cache := rediscache.New(redis)

	srv := newServer(storage, cache)

	httpPort := ":" + os.Getenv("HTTP_PORT")
	return http.ListenAndServe(httpPort, srv.router)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newRedis(cacheAddr string) *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr: cacheAddr,
	})

	return rc
}

// handlers ...

func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.ID = uuid.NewString()

	err = s.storage.User().CreateUser(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &CreateUserResponse{
		ID:   u.ID,
		Name: u.Name,
		Age:  u.Age,
		Job:  u.Job,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	u, err := s.cache.User().GetUser(id)
	// Если нашли в Кэше
	if err == nil {
		s.respondUser(w, u)
		return
	}

	// Ошибка, не связанная с отсутствием юзера в Кэше
	if !errors.Is(err, redis.Nil) {
		http.Error(w, "cache: get user error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если не нашли в Кэше
	// Пробуем искать в БД
	u, err = s.storage.User().GetUser(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}
		http.Error(w, "postgres: get user error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если нашли в БД - создаем в Кэше
	if err := s.cache.User().CreateUser(u); err != nil {
		http.Error(w, "cache: create user error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.respondUser(w, u)
}

func (s *server) respondUser(w http.ResponseWriter, u *models.User) {
	res := &GetUserResponse{
		Name: u.Name,
		Age:  u.Age,
		Job:  u.Job,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
