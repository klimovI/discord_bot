package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	hello = "hello"
	join  = "join"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        hello,
			Description: "Says hello",
		},
		{
			Name:        join,
			Description: "Joins voice channel",
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		hello: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hello!",
				},
			}

			s.InteractionRespond(i.Interaction, response)
		},
		join: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userVoiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
			if err != nil {
				log.Printf("Error getting user's voice state: %v", err)
				return
			}
			if userVoiceState == nil || userVoiceState.ChannelID == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You are not currently in a voice channel.",
					},
				})
				return
			}

			voiceConn, err := s.ChannelVoiceJoin(i.GuildID, userVoiceState.ChannelID, false, true)
			if err != nil {
				log.Printf("Error joining voice channel: %v", err)
				return
			}

			defer voiceConn.Close()

			var channelName string
			channel, err := s.Channel(userVoiceState.ChannelID)
			if err == nil {
				channelName = channel.Name
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Joined %s.", channelName),
				},
			})
		},
	}
)
