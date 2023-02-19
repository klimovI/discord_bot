package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/klimovI/discord_bot/pkg/commands"
	"github.com/klimovI/discord_bot/pkg/config"
)

type Bot struct {
	session *discordgo.Session

	commands         []*discordgo.ApplicationCommand
	commandsHandlers map[string]commands.CommandHandler
}

func main() {
	bot := newBot()
	defer bot.cleanup()

	bot.addCommands()

	log.Print("Bot started")
	log.Print("Press Ctrl+C to stop")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Print("Bot stopped")
	os.Exit(0)
}

func newBot() *Bot {
	bot := new(Bot)

	session, err := discordgo.New("Bot " + config.Config.Bot.Token)
	if err != nil {
		log.Panicf("Invalid bot parameters: %v", err)
	}

	if err := session.Open(); err != nil {
		log.Panicf("Cannot open the session: %v", err)
	}

	bot.session = session
	bot.commands = make([]*discordgo.ApplicationCommand, len(commands.Commands))
	bot.commandsHandlers = make(map[string]commands.CommandHandler, len(commands.Commands))

	return bot
}

func (bot *Bot) addCommands() {
	appID := bot.appID()

	for i, command := range commands.Commands {
		name := command.Data.Name

		createdCommand, err := bot.session.ApplicationCommandCreate(
			appID, "",
			(*discordgo.ApplicationCommand)(&command.Data),
		)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", name, err)
		}

		bot.commands[i] = createdCommand
		bot.commandsHandlers[name] = command.Handler
	}

	bot.addCommandsHandler()
}

func (bot *Bot) addCommandsHandler() {
	bot.session.AddHandler(
		func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			if interaction.Type != discordgo.InteractionApplicationCommand {
				return
			}

			name := interaction.ApplicationCommandData().Name
			handler, ok := bot.commandsHandlers[name]

			ctx := commands.Context{
				Session:     session,
				Interaction: interaction.Interaction,
			}

			if !ok {
				ctx.Respond("Command no found")
				return
			}

			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
					ctx.RespondWithError()
				}
			}()

			handler(ctx)
		},
	)
}

func (bot *Bot) cleanup() {
	bot.closeSession()
	bot.disconnectVoiceConnections()
	// b.deleteCommands()
}

func (bot *Bot) closeSession() {
	if err := bot.session.Close(); err != nil {
		log.Printf("Error closing session: %v", err)
	}
}

func (bot *Bot) disconnectVoiceConnections() {
	for _, voiceConnection := range bot.session.VoiceConnections {
		if err := voiceConnection.Disconnect(); err != nil {
			guildID := voiceConnection.GuildID
			log.Printf("Error disconnecting voice connection GuildID = '%v': %v", guildID, err)
		}
	}
}

func (bot *Bot) deleteCommands() {
	appID := bot.appID()

	for _, command := range bot.commands {
		err := bot.session.ApplicationCommandDelete(appID, "", command.ID)

		if err != nil {
			log.Printf("Error deleting '%v' command: %v", command.Name, err)
		}
	}
}

func (bot *Bot) appID() string {
	return bot.session.State.User.ID
}
