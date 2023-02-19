package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Session     *discordgo.Session
	Interaction *discordgo.Interaction
}

type CommandData discordgo.ApplicationCommand
type CommandHandler func(ctx Context)

type Command struct {
	Data    CommandData
	Handler func(ctx Context)
}

var Commands = []Command{emojiCommand, helloCommand, joinCommand, leaveCommand, testCommand}
