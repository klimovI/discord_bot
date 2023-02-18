package commands

import "github.com/bwmarrin/discordgo"

func (ctx *Context) Respond(content string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}

	ctx.Session.InteractionRespond(ctx.Interaction, response)
}
