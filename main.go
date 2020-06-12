package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("Your Telegram Token")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateChan, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}
	for update := range updateChan {
		if update.Message == nil {
			fmt.Println("It is not message")
			continue
		}
		message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		message.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(message)
		if err != nil {
			log.Panic(err)
		}
	}
}
