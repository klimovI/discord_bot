package main

import (
	"log"
	"os"
	"os/signal"

	"discord_bot/pkg/bot"
)

func main() {
	bot := new(bot.Bot)
	bot.Start()
	Quit()
}

func Quit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Quit")

	os.Exit(0)
}
