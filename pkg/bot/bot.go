package bot

import (
	"discord_bot/pkg/config"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct{}

func (b *Bot) Start() {
	bot, err := discordgo.New("Bot " + config.Config.Bot.Token)
	if err != nil {
		log.Fatal(err)
	}

	if err := bot.Open(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := bot.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		_, err := bot.ChannelMessageSend("573482163618578433", "zxc")
		if err != nil {
			log.Printf("Failed to send message: %s", err)
		}
	}()

	log.Println("Start")
}
