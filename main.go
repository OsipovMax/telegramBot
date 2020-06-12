package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var valCurs ValCurs

// ValCurs is ...
type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

// Valute is ...
type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

func getExchangeRate() ValCurs {
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_utf8.xml")
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Panic(resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err")
	}
	v := ValCurs{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return v
}

func commandHandler(command string, arg string) string {
	switch command {
	case "exchange_rate":
		switch arg {
		case "Доллар США":
			for _, val := range valCurs.Valutes {
				if val.Name == "Доллар США" {
					return val.Value
				}
			}
		case "Евро":
			for _, val := range valCurs.Valutes {
				if val.Name == "Евро" {
					return val.Value
				}
			}
		}
	}
	return ""
}

func main() {
	valCurs = getExchangeRate()
	bot, err := tgbotapi.NewBotAPI("Your Token")
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
		if update.Message.IsCommand() {
			command := update.Message.Command()
			arg := update.Message.CommandArguments()
			res := commandHandler(command, arg)
			if res == "" {
				res = "Получение информации о такой валюте невозможно"
			}
			update.Message.Text = res
			message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, err = bot.Send(message)
			if err != nil {
				log.Panic(err)
			}
		} else {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			message.ReplyToMessageID = update.Message.MessageID
			_, err = bot.Send(message)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
