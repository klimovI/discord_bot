package bot

import (
	"log"
	"os"
	"os/signal"

	cmds "github.com/klimovI/discord_bot/pkg/commands" // TODO rename to commands
	"github.com/klimovI/discord_bot/pkg/config"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session                  *discordgo.Session
	commandsHandlers         map[string]cmds.CommandHandler
	registeredCommands       []*discordgo.ApplicationCommand
	voiceConnectionByGuildID map[string]*discordgo.VoiceConnection
}

func (b *Bot) Start() {
	b.init()

	if err := b.session.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer func() {
		if err := b.session.Close(); err != nil {
			log.Println(err)
		}
	}()

	b.addHandlers()
	b.addCommands()
	b.createCommands()

	log.Print("Bot Started")
	log.Print("Press Ctrl+C to stop")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	b.cleanup()

	log.Print("Bot Stopped")
	os.Exit(0)
}

func (b *Bot) init() {
	b.initSession()

	b.voiceConnectionByGuildID = make(map[string]*discordgo.VoiceConnection)
}

func (b *Bot) initSession() {
	session, err := discordgo.New("Bot " + config.Config.Bot.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	b.session = session
}

func (bot *Bot) addCommands() {
	commands := cmds.Commands
	commandsCount := len(commands)
	// bot.registeredCommands = make([]*discordgo.ApplicationCommand, commandsCount) TODO add
	bot.commandsHandlers = make(map[string]cmds.CommandHandler, commandsCount)

	session := bot.session
	appID := session.State.User.ID

	for i, command := range commands {
		name := command.Data.Name

		cmd, err := session.ApplicationCommandCreate(appID, "", (*discordgo.ApplicationCommand)(&command.Data))
		if err != nil {
			log.Println(i, cmd) // TODO remove
			log.Panicf("Cannot create '%v' command: %v", name, err)
		}

		// bot.registeredCommands[i] = cmd TODO add
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

			ctx := cmds.Context{
				Session:     session,
				Interaction: interaction.Interaction,
			}

			if !ok {
				ctx.Respond("Command no found")
				return
			}

			handler(ctx)
		},
	)
}

func (bot *Bot) addHandlers() {
	bot.session.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Type != discordgo.InteractionApplicationCommand {
				return
			}

			commandName := i.ApplicationCommandData().Name
			handler, ok := commandsHandlers[commandName]

			if !ok {
				bot.sendMessage(s, i, "Command no found")
				return
			}

			handler(bot, s, i)
		},
	)
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

func (b *Bot) cleanup() {
	// b.deleteCommands()
	b.leaveAllVoice()
}

func (b *Bot) leaveAllVoice() {
	for _, conn := range b.voiceConnectionByGuildID {
		conn.Disconnect()
	}
}

func (b *Bot) deleteCommands() {
	for _, cmd := range b.registeredCommands {
		if err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", cmd.ID); err != nil {
			log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
