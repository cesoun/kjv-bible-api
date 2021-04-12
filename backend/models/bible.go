package models

type ErrorMessage struct {
	Error string `json:"error"`
}

type NestedBook struct {
	Title   string `json:"title"`
	Alt     string `json:"alt"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
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
