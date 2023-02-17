package bot

import (
	"discord_bot/pkg/config"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
}

func (b *Bot) Start() {
	b.initSession()

	if err := b.session.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer func() {
		if err := b.session.Close(); err != nil {
			log.Println(err)
		}
	}()

	b.addHandlers()
	b.createCommands()

	log.Print("Bot Started")
	log.Print("Press Ctrl+C to exit")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	b.deleteCommands()

	log.Print("Bot Stopped")
	os.Exit(0)
}

func (b *Bot) initSession() {
	session, err := discordgo.New("Bot " + config.Config.Bot.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	b.session = session
}

func (b *Bot) addHandlers() {
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if handler, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		}
	})
}

func (b *Bot) createCommands() {
	b.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))

	for i, command := range commands {
		cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", command.Name, err)
		}

		b.registeredCommands[i] = cmd
	}
}

func (b *Bot) deleteCommands() {
	for _, cmd := range b.registeredCommands {
		if err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", cmd.ID); err != nil {
			log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
