package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	CHANNEL_SECRET := os.Getenv("CHANNEL_SECRET")
	CHANNEL_TOKEN := os.Getenv("CHANNEL_TOKEN")

	API_base_URL := "http://ec2-3-23-60-80.us-east-2.compute.amazonaws.com/notes"

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

					// Post text to note api
					note := map[string]string{"userId": event.Source.UserID, "content": message.Text}
					json_data, err := json.Marshal(note)
					if err != nil {
						log.Fatal(err)
					}

					resp, err := http.Post(API_base_URL, "application/json", bytes.NewBuffer(json_data))
					if err != nil {
						log.Fatal(err)
					}
					defer resp.Body.Close()
					fmt.Println(resp.Body)

					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.UserID), linebot.NewTextMessage("Created!")).Do(); err != nil {
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
