package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var API_base_URL string = "http://ec2-3-23-60-80.us-east-2.compute.amazonaws.com/notes"

func GetNotes(event *linebot.Event, pager string) []Note {
	resp, err := http.Get(API_base_URL + "/" + event.Source.UserID + "/" + pager)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var notes []Note

	if err := json.Unmarshal(body, &notes); err != nil {
		log.Fatal(err)
	}
	return notes
}

func CreateReplyText(notes []Note) string {
	var replyText string

	for i, note := range notes {
		t, err := time.Parse(time.RFC3339, note.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		loc := time.FixedZone("Asia/Tokyo", 9*60*60)
		t = t.In(loc)

		if i == 0 {
			replyText += "🗓 " + t.Format("2006/01/02 15:04")
			replyText += "\n👉🏻 " + note.Content
		} else {
			replyText += "\n\n🗓 " + t.Format("2006/01/02 15:04")
			replyText += "\n👉🏻 " + note.Content
		}
	}
	return replyText
}

func CreateReplyWithMoreNotes(replyText string, pager string) linebot.SendingMessage {
	return linebot.NewTextMessage(replyText).WithQuickReplies(
		linebot.NewQuickReplyItems(
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("もっと見る", "もっと見る😉 "+pager)),
		))
}
