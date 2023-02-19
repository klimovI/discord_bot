package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (ctx *Context) Respond(content string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}

	ctx.Session.InteractionRespond(ctx.Interaction, response)
}

func (ctx *Context) RespondWithError() {
	ctx.Respond("Error")
}

func (ctx *Context) UserVoiceChannelId() (string, error) {
	guildID := ctx.GuildID()
	userID := ctx.UserID()

	userVoiceState, err := ctx.Session.State.VoiceState(guildID, userID)
	if err != nil {
		return "", fmt.Errorf("error getting use voice state: %w", err)
	}

	return userVoiceState.ChannelID, nil
}

func (ctx *Context) channelName(channelId string) (*string, error) {
	channel, err := ctx.Session.Channel(channelId)

	if err != nil {
		return nil, fmt.Errorf("error getting channel: %w", err)
	}

	return &channel.Name, nil
}

func (ctx *Context) GuildID() string {
	return ctx.Interaction.GuildID
}

func (ctx *Context) UserID() string {
	return ctx.Interaction.Member.User.ID
}

func (ctx *Context) CommandName() string {
	return ctx.Interaction.ApplicationCommandData().Name
}

func (ctx *Context) joinVoiceChannel(guildID string, channelID string) (*discordgo.VoiceConnection, error) {
	voiceConnection, err := ctx.Session.ChannelVoiceJoin(guildID, channelID, false, true)

	if err != nil {
		return nil, fmt.Errorf("error joining voice channel: %w", err)
	}

	ctx.Session.VoiceConnections[guildID] = voiceConnection

	return voiceConnection, nil
}
