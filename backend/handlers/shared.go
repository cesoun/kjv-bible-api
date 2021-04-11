package handlers

import (
	"fmt"
	"math/rand"
	"strings"
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

func GetFirstWordOccurrence(data *models.BookData, word string) *models.NestedBook {
	// Get the books.
	books := gjson.Get(data.Json, `books`)
	for _, book := range books.Array() {
		// Store title & alt
		title := gjson.Get(book.Raw, `title`)
		alt := gjson.Get(book.Raw, `alt`)

		// Get the chapters.
		chapters := gjson.Get(book.Raw, `chapters`)
		for _, chapter := range chapters.Array() {
			// Store chapter #
			c := gjson.Get(chapter.Raw, `chapter`)

			// Get the verses.
			verses := gjson.Get(chapter.Raw, `verses`)
			for v, verse := range verses.Array() {
				// See if the verse contains the word.
				words := strings.Fields(strings.ToLower(verse.Str))
				for _, w := range words {
					// Found a match? return struct.
					if word == w {
						return &models.NestedBook{
							Title:   title.Str,
							Alt:     alt.Str,
							Chapter: int(c.Int()),
							Verse:   v + 1,
						}
					}
				}
			}
		}
	}

	// Didn't find anything, nil.
	return nil
}
