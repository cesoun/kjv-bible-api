package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cesoun/kjv-bible-api/handlers"
	"github.com/cesoun/kjv-bible-api/models"
	"github.com/gorilla/mux"
)

type App struct {
	Router  *mux.Router
	Server  *http.Server
	addr    string
	port    string
	newTest *models.BookData
	oldTest *models.BookData
}

func (a *App) Init(addr, port string) {
	a.Router = mux.NewRouter()

	a.addr = addr
	a.port = port

	a.loadJson()
	a.initRoutes()

	a.Server = &http.Server{
		Handler:      a.Router,
		Addr:         fmt.Sprintf("%s:%s", addr, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (a *App) Run() {
	log.Fatal(a.Server.ListenAndServe())
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getBase(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *App) getRandomVerse(w http.ResponseWriter, r *http.Request) {
	// seed random and pick a version.
	rand.Seed(time.Now().UnixNano())
	rng := rand.Intn(2)

	var verse models.RandomVerse

	// Determine which testament to take from.
	if rng == 0 {
		verse = handlers.GetRandomVerse(a.oldTest)
	} else {
		verse = handlers.GetRandomVerse(a.newTest)
	}

	respondWithJSON(w, http.StatusOK, verse)
}

func (a *App) loadJson() {
	// Load the New Testament JSON
	r, err := os.ReadFile("../data/kjv-new-min.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup and assign the data.
	newData := &models.BookData{
		Json:    string(r),
		Version: "The New Testament of the King James Bible",
	}
	a.newTest = newData

	// Load the Old Testament JSON
	r, err = os.ReadFile("../data/kjv-old-min.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup and assign the data.
	oldData := &models.BookData{
		Json:    string(r),
		Version: "The Old Testament of the King James Version of the Bible",
	}
	a.oldTest = oldData
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", a.getBase).Methods("GET")

	api := a.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", a.getBase).Methods("GET")
	api.HandleFunc("/random", a.getRandomVerse).Methods("GET")
}
