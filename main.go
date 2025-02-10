package main

import (
	"log"
	"github.com/bwmarrin/discordgo"
	"os"
)

func main() {
	initDatabase()
	bootstrap()
	go startWikipediaStream()

	token := os.Getenv("DISCORD_BOT_TOKEN")

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is online as", s.State.User.Username)
	})

	discord.AddHandler(commandHandler)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}

	log.Println("Bot is running...")
	select {}
}
