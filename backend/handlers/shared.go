package handlers

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/cesoun/kjv-bible-api/models"
	"github.com/tidwall/gjson"
)

func GetRandomVerse(data *models.BookData) models.BibleVerse {
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

	return models.BibleVerse{
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
	// Match a whole word and regex is nicer for this.
	re := regexp.MustCompile(fmt.Sprintf(`(\b%s\b)`, word))

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

				// Regular expression matching is probably better here.
				if re.MatchString(verse.Str) {
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

	// Didn't find anything, nil.
	return nil
}

func GetVerse(data *models.BookData, book string, chapter int, verse int) (*models.BibleVerse, error) {
	// Reject 0 chapter/verse.
	if chapter < 1 {
		return nil, errors.New(fmt.Sprintf("Invalid chapter: %d", chapter))
	}

	if verse < 1 {
		return nil, errors.New(fmt.Sprintf("Invalid verse: %d", verse))
	}

	// Find the book.
	titles := gjson.Get(data.Json, `titles`)
	for _, titleSet := range titles.Array() {
		title := gjson.Get(titleSet.Raw, `title`)
		alt := gjson.Get(titleSet.Raw, `alt`)

		// Does the book exist?
		if stringsAreEqual(book, title.Str) || stringsAreEqual(book, alt.Str) {
			// Collect the book.
			b := gjson.Get(data.Json, fmt.Sprintf(`books.#(title=="%s")`, title.Str))
			if !b.Exists() {
				b = gjson.Get(data.Json, fmt.Sprintf(`books.#(title=="%s")`, alt.Str))
			}

			// Collect the chapters & determine if we can access the desired chapter.
			c := gjson.Get(b.Raw, `chapters.#`)
			if chapter > int(c.Int()) {
				return nil, errors.New(fmt.Sprintf("Could not find chapter: %d in book: %s", chapter, book))
			}

			// We do store chapter # in the json, however it's probably faster to just direct access this.
			c = gjson.Get(b.Raw, fmt.Sprintf(`chapters.%d`, chapter-1))

			// Collect the verses & determine if we can access the desired verse.
			v := gjson.Get(c.Raw, `verses.#`)
			if verse > int(v.Int()) {
				return nil, errors.New(fmt.Sprintf("Could not find verse: %d in book: %s, chapter: %d", verse, book, chapter))
			}

			// Get the verse.
			v = gjson.Get(c.Raw, fmt.Sprintf(`verses.%d`, verse-1))

			// Return that bad boy.
			return &models.BibleVerse{
				Version: data.Version,
				Book: models.NestedBook{
					Title:   title.Str,
					Alt:     alt.Str,
					Chapter: chapter,
					Verse:   verse,
				},
				Verse: v.Str,
			}, nil
		}
	}

	// Didn't find the given book, return.
	return nil, errors.New(fmt.Sprintf("Could not find book: %s", book))
}

func stringsAreEqual(a string, b string) bool {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == 0
}
