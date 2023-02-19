package commands

import "log"

var leaveCommand = Command{
	Data: CommandData{
		Name:        "leave",
		Description: "Leave voice chanel",
	},
	Handler: func(ctx Context) {
		guildID := ctx.GuildID()
		voiceConnection, ok := ctx.Session.VoiceConnections[guildID]

		if !ok {
			ctx.Respond("Not connected to a voice channel")
			return
		}

		if err := voiceConnection.Disconnect(); err != nil {
			log.Println(err)
			ctx.RespondWithError()
			return
		}

		delete(ctx.Session.VoiceConnections, guildID)

		ctx.Respond("Left voice channel")
	},
}
