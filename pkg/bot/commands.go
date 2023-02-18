package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	hello = "hello"
	join  = "join"
	leave = "leave"
	emoji = "emoji"
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
		{
			Name:        leave,
			Description: "Leave voice chanel",
		},
		{
			Name:        emoji,
			Description: "Send emoji to chat",
		},
	}
	commandsHandlers = map[string]func(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate){
		hello: helloHandler,
		join:  joinHandler,
		leave: leaveHandler,
		emoji: emojiHandler,
	}
)

func helloHandler(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.sendMessage(s, i, "Hello")
}

func joinHandler(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {

	channelID := getUserVoiceChannelId(b, s, i)
	if channelID == nil {
		return
	}

	existingConnection, ok := b.voiceConnectionByGuildID[i.GuildID]
	if ok && existingConnection.ChannelID == *channelID {
		b.sendMessage(s, i, "Already connected")
		return
	}

	voiceConn, err := s.ChannelVoiceJoin(i.GuildID, *channelID, false, true)
	if err != nil {
		log.Printf("Error joining voice channel: %v", err)
		return
	}

	b.voiceConnectionByGuildID[i.GuildID] = voiceConn

	var channelName string
	channel, err := s.Channel(*channelID)
	if err == nil {
		channelName = channel.Name
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Joined %s.", channelName),
		},
	})
}

func getUserVoiceChannelId(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) *string {
	userVoiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)

	if err != nil {
		log.Printf("Error getting use voice state: %v", err)
		return nil
	}

	if userVoiceState == nil || userVoiceState.ChannelID == "" {
		b.sendMessage(s, i, "You are not currently in a voice channel")
		return nil
	}

	return &userVoiceState.ChannelID
}

func leaveHandler(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	voiceConn, ok := b.voiceConnectionByGuildID[i.GuildID]

	if !ok {
		b.sendMessage(s, i, "Not connected")
		return
	}

	voiceConn.Disconnect()

	delete(b.voiceConnectionByGuildID, i.GuildID)

	b.sendMessage(s, i, "Left channel")
}

func emojiHandler(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.sendMessage(s, i, "🤡")
}

func (b *Bot) sendMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}

	s.InteractionRespond(i.Interaction, response)
}
