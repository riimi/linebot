// Package p contains a Pub/Sub Cloud Function.
package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
	"time"
)

const MAXITEM int = 5

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	//log.Println(string(m.Data))
	//projectID := `test`
	projectID := os.Getenv("PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer client.Close()

	repoRssService := &RssServiceRepoFirestore{Client: client}

	if err := repoRssService.Foreach(HandleService(client)); err != nil {
		log.Fatalf("%v", err)
	}

	return nil
}

func HandleService(client *firestore.Client) func(service RssService) error {
	repoRssItem := &RssItemRepoFirestore{Client: client}

	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	return func(service RssService) error {
		nItem := 0
		feed, err := GetRssFeed(service.Url)
		if err != nil {
			return err
		}

		for _, item := range feed.Items {
			newItem := RssItem{
				ServiceName: service.Name,
				Title:       item.Title,
				Link:        item.Link,
				Published:   item.Published,
				CreatedAt:   time.Now(),
			}

			if !repoRssItem.IsNewItem(newItem) || nItem >= MAXITEM {
				log.Printf("%s got %d new items", service.Name, nItem)
				return nil
			}

			if err := repoRssItem.Add(newItem); err != nil {
				return err
			}
			nItem += 1
			log.Printf("[new item %d] %v", nItem, newItem)

			for _, id := range service.UserID {
				log.Print(id)
				if _, err := bot.PushMessage(id, linebot.NewTextMessage(newItem.TextMessage())).Do(); err != nil {
					return err
				}
			}
		}

		return nil
	}
}
