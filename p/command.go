package p

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var cmdRssService = &cobra.Command{
	Use:   "Rss",
	Short: "Add/Remove/List/Subscribe/Unsubscribe Rss feed",
	Long:  `Add/Remove/List/Subscribe/Unsubscribe Rss feed`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var cmdRssAdd = &cobra.Command{
	Use:     "add [name] [url]",
	Aliases: []string{"등록", "추가"},
	Short:   "Add Rss feed",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName, serviceUrl := args[0], args[1]
		out, err := exec.Command("curl", "-c", "cookie", "-XGET", "-L", serviceUrl).Output()
		if err != nil {
			log.Print("curl:" + err.Error())
			fmt.Print(err)
			return
		}
		fp := gofeed.NewParser()
		feed, err := fp.ParseString(string(out))
		if err != nil {
			log.Print(string(out))
			log.Print(err)
			fmt.Print(err)
			return
		}
		ReplyMessage(feed.Title + "\n" + feed.Description)
		model.CreateService(model.RssService{
			Name:        serviceName,
			Url:         serviceUrl,
			Title:       feed.Title,
			Description: feed.Description,
			Link:        feed.Link,
			Language:    feed.Language,
		})
		service, _ := model.GetServiceByName(serviceName)
		model.CreateRssItems(feed, service, false)

	},
}

var cmdRssRemove = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"삭제", "제거"},
	Short:   "Remove Rss feed",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := model.DeleteServiceByName(args[0]); err == nil {
			cmd.Printf("rss feed [%s] is removed\n", args[0])
		} else {
			cmd.Printf("failed to remove the service")
		}
	},
}

var cmdRssList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"목록", "리스트"},
	Short:   "list Rss feed",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		services := model.AllServices()
		//raw, _ := json.Marshal(services)
		//ReplyMessage(string(raw))
		/*
			cmd.Println("# rss list")
			for _, s := range services {
				cmd.Println(s)
			}
		*/
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("서버 에러!", template.FlexContainerRssServices(services, SourceToID())),
		).Do(); err != nil {
			log.Print(err)
		}
	},
}

var cmdRssSubscribe = &cobra.Command{
	Use:     "subscribe [name]",
	Aliases: []string{"구독"},
	Short:   "Subscribe Rss feed",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		service, err := model.GetServiceByName(serviceName)
		if err != nil {
			cmd.Printf("the service %s is not available\n", serviceName)
			return
		}
		id := SourceToID()
		if err := model.Subscribe(id, service.Name); err != nil {
			log.Print(err)
			cmd.Printf("failed to subscribe")
		} else {
			cmd.Printf("subscribe %s\n", serviceName)
		}
	},
}

var cmdRssUnsubscribe = &cobra.Command{
	Use:     "unsubscribe [name]",
	Aliases: []string{"구독해제"},
	Short:   "Unsubscribe Rss feed",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		service, err := model.GetServiceByName(serviceName)
		if err != nil {
			cmd.Printf("the service %s is not available\n", serviceName)
			return
		}
		id := SourceToID()
		if err := model.Unsubscribe(id, service.Name); err != nil {
			cmd.Print("failed to unsubscribe")
		} else {
			cmd.Printf("unsubscribe %s\n", serviceName)
		}
	},
}
