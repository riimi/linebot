package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"p/consumer"
)

func main() {
	projectID := `test`

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer client.Close()

	//repoRssService := &consumer.RssServiceRepoFirestore{Client: client}
	//repoRssItem := &consumer.RssItemRepoFirestore{Client: client}
	repoSubs := &consumer.SubscriptionRepoFirestore{Client: client}

	repoSubs.Add(consumer.Subscription{
		UserId:      []string{},
		ServiceName: `horriblesubs-1080p`,
	})

	//cases := map[string]string{
	//	`clien-hotdeal`: `http://feeds.feedburner.com/c_shop`,
	//	`horriblesubs-1080p`: `http://www.horriblesubs.info/rss.php?res=1080`,
	//	`ruliweb-hotdeal`: `http://bbs.ruliweb.com/market/board/1020/rss`,
	//}
	//
	//for name, url := range cases {
	//	if err := consumer.AddRssService(repoRssService, name, url); err != nil {
	//		log.Fatal(err)
	//	}
	//}

	//if err := repoRssService.Foreach(func(service consumer.RssService) error {
	//	log.Printf("%v", service)
	//
	//	//url := `http://bbs.ruliweb.com/market/board/1020/rss`
	//	feed, err := consumer.GetRssFeed(service.Url)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, item := range feed.Items {
	//		newItem := consumer.RssItem{
	//			ServiceName: service.Name,
	//			Title:       item.Title,
	//			Link:        item.Link,
	//			Published:   item.Published,
	//			CreatedAt:   time.Now(),
	//		}
	//
	//		_, err := repoRssItem.Get(newItem.Link)
	//		if err == nil {
	//			return nil
	//		}
	//		stat, ok := status.FromError(err)
	//		if !ok || stat.Code() != codes.NotFound {
	//			return stat.Err()
	//		}
	//
	//		if err := repoRssItem.Add(newItem); err != nil {
	//			return err
	//		}
	//
	//		sub, err := repoSubs.Get(service.Name)
	//		if  err != nil {
	//			return err
	//		}
	//		for _, id := range sub.UserId {
	//			log.Print(id)
	//		}
	//		break
	//	}
	//
	//	return nil
	//}); err != nil {
	//	log.Fatalf("%v", err)
	//}
}
