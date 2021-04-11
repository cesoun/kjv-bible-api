package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cesoun/kjv-bible-api/models"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

type BookData struct {
	rawJson string
	version string
}

type App struct {
	Router  *mux.Router
	Server  *http.Server
	addr    string
	port    string
	newTest *BookData
	oldTest *BookData
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
	// type result struct {}

	// seed random and pick a version.
	rand.Seed(time.Now().UnixNano())
	rng := rand.Intn(2)

	var version *BookData

	// Determine which testament to take from.
	if rng == 0 {
		version = a.oldTest
	} else {
		version = a.newTest
	}

	// Get the number of books in version & pick a random one.
	books := gjson.Get(version.rawJson, `books.#`)
	book := rand.Intn(int(books.Int()))

	// Get the json for random book
	bookJson := gjson.Get(version.rawJson, fmt.Sprintf("books.%d", book))

	// Get title & alt title
	title := gjson.Get(bookJson.Raw, "title")
	alt := gjson.Get(bookJson.Raw, "alt")

	// Pick random Chatper
	chapters := gjson.Get(bookJson.Raw, "chapters.#")
	chapter := rand.Intn(int(chapters.Int()))

	// Pick a random verse.
	verses := gjson.Get(bookJson.Raw, fmt.Sprintf("chapters.%d.verses.#", chapter))
	verse := rand.Intn(int(verses.Int()))

	randVerse := gjson.Get(bookJson.Raw, fmt.Sprintf("chapters.%d.verses.%d", chapter, verse))

	respondWithJSON(w, http.StatusOK, models.RandomVerse{
		Version: *&version.version,
		Book: models.NestedBook{
			Title:   title.Str,
			Alt:     alt.Str,
			Chapter: chapter + 1,
			Verse:   verse + 1,
		},
		Verse: randVerse.Str,
	})
}

func (a *App) loadJson() {
	// Load the New Testament JSON
	r, err := os.ReadFile("../data/kjv-new-min.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup and assign the data.
	newData := &BookData{
		rawJson: string(r),
		version: "The New Testament of the King James Bible",
	}
	a.newTest = newData

	// Load the Old Testament JSON
	r, err = os.ReadFile("../data/kjv-old-min.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup and assign the data.
	oldData := &BookData{
		rawJson: string(r),
		version: "The Old Testament of the King James Version of the Bible",
	}
	a.oldTest = oldData
}

func (a *App) initRoutes() {
	// TODO: Might breakout handlers to shrink this file down.

	// /
	a.Router.HandleFunc("/", a.getBase).Methods("GET")

	// /api
	api := a.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/random", a.getRandomVerse).Methods("GET")
}
