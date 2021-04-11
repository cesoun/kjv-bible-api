package handlers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cesoun/kjv-bible-api/models"
	"github.com/tidwall/gjson"
)

func GetRandomVerse(data *models.BookData) models.RandomVerse {
	// Seed rand before starting.
	rand.Seed(time.Now().UnixNano())

	// Get the number of books in version & pick a random one.
	books := gjson.Get(data.Json, `books.#`)
	book := rand.Intn(int(books.Int()))

	// Get the json for random book
	bookJson := gjson.Get(data.Json, fmt.Sprintf("books.%d", book))

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

	return models.RandomVerse{
		Version: data.Version,
		Book: models.NestedBook{
			Title:   title.Str,
			Alt:     alt.Str,
			Chapter: chapter + 1,
			Verse:   verse + 1,
		},
		Verse: randVerse.Str,
	}
}
