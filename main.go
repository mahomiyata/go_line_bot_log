package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Note struct {
	UserID    string `json:"UserId"`
	Content   string `json:"Content"`
	CreatedAt string `json:"CreatedAt"`
}

func main() {
	CHANNEL_SECRET := os.Getenv("CHANNEL_SECRET")
	CHANNEL_TOKEN := os.Getenv("CHANNEL_TOKEN")

	bot, err := linebot.New(CHANNEL_SECRET, CHANNEL_TOKEN)

	if err != nil {
		log.Fatal(err)
	}

	isExecutedTime := flag.Bool("broadcast", false, "a bool for checking broadcast")
	flag.Parse()

	if *isExecutedTime {
		feelingLog := linebot.NewTextMessage("今はどんな感じ？").WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton("", linebot.NewMessageAction("良い感じ🐥", "良い感じ🐥")),
				linebot.NewQuickReplyButton("", linebot.NewMessageAction("まあまあ🐣", "まあまあ🐣")),
				linebot.NewQuickReplyButton("", linebot.NewMessageAction("だめかな……🐤", "だめかな……🐤")),
			))
		if _, err := bot.BroadcastMessage(feelingLog).Do(); err != nil {
			log.Fatal()
		}
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
					// Reply 5 latest notes
					if strings.Contains(message.Text, "★履歴★") {
						var notes []Note = GetNotes(event, "1")
						var replyText string = CreateReplyText(notes)

						ReplyMessage := CreateReplyWithMoreNotes(replyText, "2")
						if _, err := bot.ReplyMessage(event.ReplyToken, ReplyMessage).Do(); err != nil {
							log.Fatal(err)
						}

						// Reply more notes
					} else if strings.Contains(message.Text, "もっと見る😉") {
						splitStr := strings.Split(message.Text, " ")

						var notes []Note = GetNotes(event, splitStr[1])

						// FIXME: You can use case!
						if len(notes) == 0 {
							if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("もう無いみたい🐬")).Do(); err != nil {
								log.Print(err)
							}
						} else {
							var replyText string = CreateReplyText(notes)

							next, err := strconv.Atoi(splitStr[1])
							if err != nil {
								log.Fatal(err)
							}
							next += 1

							replyMessage := CreateReplyWithMoreNotes(replyText, strconv.Itoa(next))

							if _, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
								log.Fatal(err)
							}
						}

					} else {
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

						if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("789", "10863")).Do(); err != nil {
							log.Print(err)
						}
					}
				case *linebot.StickerMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("789", "10877")).Do(); err != nil {
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
