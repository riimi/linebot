package p

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"os"
)

func main() {
	//scan := bufio.NewScanner(os.Stdin)
	//for scan.Scan() {
	//	text := scan.Text()
	//	args := strings.Fields(text)
	//	if _, _, err := rootCmd.Find(args); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	rootCmd.SetArgs(args)
	//	if err := rootCmd.Execute(); err != nil {
	//		log.Fatal(err)
	//	}
	//}

	//repoRssService := &consumer.RssServiceRepoFirestore{Client: client}
	//repoRssItem := &consumer.RssItemRepoFirestore{Client: client}

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

	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer client.Close()
	//
	repoRssService := &RssServiceRepoFirestore{Client: client}
	//service, err := repoRssService.Get(`horriblesubs-1080p`)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//service.AddSubscriber(`U80c288156ed296cfa61e8325df0e271c`)
	//if err := repoRssService.Add(service); err != nil {
	//	log.Fatal(err)
	//}

	if err := repoRssService.Foreach(HandleService(client)); err != nil {
		log.Fatalf("%v", err)
	}
}
