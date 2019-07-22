package p

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := bot.ParseRequest(r)
	if err == linebot.ErrInvalidSignature {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch event.Message.(type) {
			case *linebot.TextMessage:
				if err := HandleTextMessage(bot, event); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func HandleTextMessage(bot *linebot.Client, e *linebot.Event) error {
	msg := e.Message.(*linebot.TextMessage)
	_, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(msg.Text)).Do()
	return err
}
