package commands

import (
	"fmt"
	"log"
)

var joinCommand = Command{
	Data: CommandData{
		Name:        "join",
		Description: "Joins voice channel",
	},
	Handler: func(ctx Context) {
		channelID, err := ctx.UserVoiceChannelId()
		if err != nil {
			log.Println(err)
			ctx.Respond("You are not currently in a voice channel")
			return
		}

		guildID := ctx.GuildID()
		currVoiceConn, ok := ctx.Session.VoiceConnections[guildID]
		if ok && currVoiceConn.ChannelID == channelID {
			ctx.Respond("Already connected")
			return
		}

		if _, err := ctx.joinVoiceChannel(guildID, channelID); err != nil {
			log.Println(err)
			ctx.RespondWithError()
			return
		}

		channelName, err := ctx.channelName(channelID)
		if err != nil {
			log.Println(err)
			ctx.RespondWithError()
			return
		}

		response := fmt.Sprintf("Joined %s.", *channelName)
		ctx.Respond(response)
	},
}
