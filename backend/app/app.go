package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	a.Router.StrictSlash(true)

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

func (a *App) loadJson() {
	// Load the New Testament JSON
	r, err := os.ReadFile("./data/kjv-new-min.json")
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
	r, err = os.ReadFile("./data/kjv-old-min.json")
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
	// TODO: Caching/Storing these results would be ideal. They aren't going to change unless we change the data.
	api.HandleFunc("/firstoccur/{word:[a-zA-Z]+}", a.getFirstOccurrence).Methods("GET")
	api.HandleFunc("/totaloccur/{word:[a-zA-Z]+}", a.getTotalOccurrence).Methods("GET")

	oldTest := api.PathPrefix("/old").Subrouter()
	oldTest.HandleFunc("/verse/{book:[a-zA-Z \\d]+}/{chapter:[\\d]+}/{verse:[\\d]+}", a.getVerse).Methods("GET")

	newTest := api.PathPrefix("/new").Subrouter()
	newTest.HandleFunc("/verse/{book:[a-zA-Z \\d]+}/{chapter:[\\d]+}/{verse:[\\d]+}", a.getVerse).Methods("GET")
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

	var verse models.BibleVerse

	// Determine which testament to take from.
	if rng == 0 {
		verse = handlers.GetRandomVerse(a.oldTest)
	} else {
		verse = handlers.GetRandomVerse(a.newTest)
	}

	respondWithJSON(w, http.StatusOK, verse)
}

func (a *App) getFirstOccurrence(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	word := params["word"]

	payload := models.FirstOccurrence{}

	// Get the old & new testament first occurrence.
	payload.OldBook = handlers.GetFirstWordOccurrence(a.oldTest, word)
	payload.NewBook = handlers.GetFirstWordOccurrence(a.newTest, word)

	// If we don't have results for either, just 204 no content.
	if payload.OldBook == nil && payload.NewBook == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Return the object with first occurrence.
	respondWithJSON(w, http.StatusOK, payload)
}

func (a *App) getTotalOccurrence(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	word := params["word"]

	payload := models.TotalOccurrence{}

	// Get the old & new testament total occurrence.
	payload.OldBook = handlers.GetWordTotalOccurrence(a.oldTest, word)
	payload.NewBook = handlers.GetWordTotalOccurrence(a.newTest, word)

	// No payload, 204.
	if payload.OldBook == nil && payload.NewBook == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Both payloads? Add a combined field.
	if payload.OldBook != nil && payload.NewBook != nil {
		payload.Both = &models.WordCount{
			Count: payload.OldBook.Count + payload.NewBook.Count,
			Word:  word,
		}
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (a *App) getVerse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	path := r.URL.Path

	// Collect params
	book := params["book"]
	chapter := params["chapter"]
	verse := params["verse"]

	// Validate chapter
	c, err := strconv.Atoi(chapter)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, models.ErrorMessage{Error: "Failed to parse chapter."})
		return
	}

	// Validate verse
	v, err := strconv.Atoi(verse)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, models.ErrorMessage{Error: "Failed to parse verse."})
		return
	}

	var payload *models.BibleVerse

	// Determine what book to look in.
	if strings.Contains(path, "/old/") {
		payload, err = handlers.GetVerse(a.oldTest, book, c, v)
	} else {
		payload, err = handlers.GetVerse(a.newTest, book, c, v)
	}

	// Respond with error if !nil
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, models.ErrorMessage{Error: err.Error()})
		return
	}

	// Respond with payload otherwise.
	respondWithJSON(w, http.StatusOK, payload)
}
