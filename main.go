package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/klimovI/discord_bot/pkg/commands"
	"github.com/klimovI/discord_bot/pkg/config"
	"github.com/klimovI/discord_bot/pkg/logger"
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

	logger.Info("Bot started")
	logger.Info("Press Ctrl+C to stop")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	logger.Info("Bot stopped")
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

			ctx := commands.Context{
				Session:     session,
				Interaction: interaction.Interaction,
			}

			commandName := ctx.CommandName()
			handler, ok := bot.commandsHandlers[commandName]

			if !ok {
				ctx.Respond("Command no found")
				return
			}

			defer func() {
				if err := recover(); err != nil {
					logger.Error(err)
					ctx.RespondWithError()
				}
			}()

			zero := 0
			print(1 / zero)

			handler(ctx)
		},
	)
}

func (bot *Bot) cleanup() {
	bot.disconnectVoiceConnections()
	// TODO uncomment
	// b.deleteCommands()
	bot.closeSession()
}

func (bot *Bot) closeSession() {
	if err := bot.session.Close(); err != nil {
		logger.Errorf("Error closing session: %v\n", err)
	}
}

func (bot *Bot) disconnectVoiceConnections() {
	for _, voiceConnection := range bot.session.VoiceConnections {
		if err := voiceConnection.Disconnect(); err != nil {
			guildID := voiceConnection.GuildID
			logger.Errorf("Error disconnecting voice connection GuildID = '%v': %v\n", guildID, err)
		}
	}
}

func (bot *Bot) deleteCommands() {
	appID := bot.appID()

	for _, command := range bot.commands {
		err := bot.session.ApplicationCommandDelete(appID, "", command.ID)

		if err != nil {
			logger.Errorf("Error deleting '%v' command: %v\n", command.Name, err)
		}
	}
}

func (bot *Bot) appID() string {
	return bot.session.State.User.ID
}
