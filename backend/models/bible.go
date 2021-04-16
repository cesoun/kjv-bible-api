package models

type ErrorMessage struct {
	Error string `json:"error"`
}

type NestedBook struct {
	Title    string `json:"title"`
	Alt      string `json:"alt"`
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse"`
	VerseUrl string `json:"verse_url,omitempty"`
}

type BibleVerse struct {
	Version string     `json:"version"`
	Book    NestedBook `json:"book"`
	Verse   string     `json:"verse"`
}

type FirstOccurrence struct {
	OldBook *NestedBook `json:"old,omitempty"`
	NewBook *NestedBook `json:"new,omitempty"`
}

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type TotalOccurrence struct {
	OldBook *WordCount `json:"old,omitempty"`
	NewBook *WordCount `json:"new,omitempty"`
	Both    *WordCount `json:"both,omitempty"`
}
