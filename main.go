package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	CHANNEL_SECRET := os.Getenv("CHANNEL_SECRET")
	CHANNEL_TOKEN := os.Getenv("CHANNEL_TOKEN")

	bot, err := linebot.New(CHANNEL_SECRET, CHANNEL_TOKEN)

	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP server for webhook
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.UserID), linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
