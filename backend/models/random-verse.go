package models

type NestedBook struct {
	Title   string `json:"title"`
	Alt     string `json:"alt"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
}

type RandomVerse struct {
	Version string     `json:"version"`
	Book    NestedBook `json:"book"`
	Verse   string     `json:"verse"`
}
