package main

type Stat struct {
    Date   string `json:"date"`
    Lang   string `json:"lang"`
    Edits  int    `json:"edits"`
    Offset int    `json:"offset"`
}

type WikipediaChange struct {
	Title string `json:"title"`
	URL   string `json:"meta.uri"`
	User  string `json:"user"`
	Meta  struct {
		Domain string `json:"domain"`
		Offset int `json:"offset"`
	} `json:"meta"`
	Time int `json:"timestamp"`
}