package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"sync"
	"github.com/bwmarrin/discordgo"
)

var changeLock = sync.Mutex{}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "!recent":
		lang := getUserLanguage(m.Author.ID)
		if len(args) > 1 {
			lang = args[1]
		}

		changeLock.Lock()
		if len(recentChanges[lang]) == 0 {
			log.Println("No changes in cache, forcing update...")
			startWikipediaStream()
		}
		changes, exists := recentChanges[lang]
		changeLock.Unlock()

		if !exists || len(changes) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No recent changes found for language: "+lang)
			return
		}

		response := "**Recent Wikipedia Edits (" + lang + ")**\n"
		limit := RECENT_CHANGE_SIZE
		if len(changes) < limit {
			limit = len(changes)
		}
		for _, change := range changes[:limit] {
			timestamp := time.Unix(int64(change.Time), 0)
			formattedTime := timestamp.Format("2006-01-02 15:04:05")
			domain := change.Meta.Domain
			if change.URL == "" {
				change.URL = "https://" + domain + "/wiki/" + change.Title
			}
			response += fmt.Sprintf("**%s** by `%s`\n[View Change](%s)\nTime: %s\n\n", change.Title, change.User, change.URL, formattedTime)
		}
		s.ChannelMessageSend(m.ChannelID, response)

	case "!setLang":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: `!setLang [language_code]`")
			return
		}
		setUserLanguage(m.Author.ID, args[1])
		s.ChannelMessageSend(m.ChannelID, "Your language has been set to: "+args[1])

	case "!stats":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: `!stats [yyyy-mm-dd]`")
			return
		}
		lang := getUserLanguage(m.Author.ID)
		count := getEditCount(lang, args[1])
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Wikipedia edits on %s (%s): **%d**", args[1], lang, count))
	}
}
