package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"github.com/r3labs/sse/v2"
)
// Constants
var MAX_BATCH_SIZE = 100
var DEFAULT_LANGUAGE = "en"
var RECENT_CHANGE_SIZE = 5

var recentChanges = make(map[string][]WikipediaChange) // Stores latest changes per language
var changeCount = make(map[string]int)
var changeOffset = make(map[string]int)

func startWikipediaStream() {
	client := sse.NewClient("https://stream.wikimedia.org/v2/stream/recentchange")

	client.SubscribeRaw(func(msg *sse.Event) {
		var change WikipediaChange
		if err := json.Unmarshal(msg.Data, &change); err != nil {
			log.Println("Error unmarshalling data:", err)
			return
		}

		change.Title = strings.ReplaceAll(change.Title, " ", "_") // This is so that links do not break because of spaces
		lang := strings.Split(change.Meta.Domain, ".")[0]

		changeLock.Lock()
		recentChanges[lang] = append([]WikipediaChange{change}, recentChanges[lang]...)
		if len(recentChanges[lang]) > RECENT_CHANGE_SIZE {
			recentChanges[lang] = recentChanges[lang][:RECENT_CHANGE_SIZE]
		}
		changeLock.Unlock()

		date := time.Unix(int64(change.Time), 0).UTC().Format("2006-01-02")

		key := makeKey(lang, date)

		changeCount[key] = changeCount[key] + 1

		if (changeCount[key] % MAX_BATCH_SIZE == 0) {
			storeEditCounts(lang, date, change.Meta.Offset)
		}
	})
}

// Bootstraps offset and count for each language-date pair from db.
func bootstrap() {
	bootstrapStats, err := getAllStats()
	if (err != nil) {
		log.Println("Error while bootstrapping", err)
		return
	}

	for _, stat := range bootstrapStats {
		key := makeKey(stat.Lang, stat.Date)
		changeCount[key] = stat.Edits
		changeOffset[key] = stat.Offset
	}
	log.Println("Offset and Edit Count maps are initialized from db!")
}

func makeKey(lang string, date string) string {
	return lang + ":" + date
}

func getEditCount(lang string, date string) int {
	return changeCount[makeKey(lang, date)]
}
