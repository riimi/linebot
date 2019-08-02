package p

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strings"
)

type Context struct {
	Linebot   *linebot.Client
	LineEvent *linebot.Event
	Firestore *firestore.Client
	UserID    string
}

func NewContext() *Context {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Fatalf("%v", err)
	}

	projectID := os.Getenv("PROJECT_ID")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &Context{
		Linebot:   bot,
		Firestore: client,
	}
}

var Ctx *Context

func init() {
	Ctx = NewContext()
}

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	events, err := Ctx.Linebot.ParseRequest(r)
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
				if err := HandleTextMessage(Ctx.Linebot, event, event.Message.(*linebot.TextMessage).Text); err != nil {
					log.Print(err)
				}
			}
		} else if event.Type == linebot.EventTypePostback {
			if err := HandleTextMessage(Ctx.Linebot, event, event.Postback.Data); err != nil {
				log.Print(err)
			}
		}
	}
}

func HandleTextMessage(bot *linebot.Client, e *linebot.Event, msg string) error {
	log.Printf("%#v\n", *e.Source)
	Ctx.UserID = SourceID(e.Source)

	Ctx.LineEvent = e
	log.Print(msg)
	args := strings.Fields(msg)
	if _, _, err := rootCmd.Find(args); err != nil {
		return err
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	if buf.Len() > 0 {
		if _, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(buf.String())).Do(); err != nil {
			return err
		}
	}

	return nil
}
