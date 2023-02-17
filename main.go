package main

import (
	"discord_bot/pkg/bot"
)

func main() {
	bot := new(bot.Bot)
	bot.Start()
}
