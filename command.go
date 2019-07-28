package p

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	cmdRssService.AddCommand(cmdRssAdd, cmdRssList, cmdRssRemove, cmdRssSubscribe, cmdRssUnsubscribe)
	rootCmd.AddCommand(cmdRssService)
}

var rootCmd = &cobra.Command{
	Use: "bot",
}

var cmdRssService = &cobra.Command{
	Use:   "!rss",
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
		if err := AddRssService(Ctx, serviceName, serviceUrl); err != nil {
			log.Print(err)
			Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("failed to add")).Do()
			return
		}
		Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("success")).Do()
	},
}

var cmdRssRemove = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"삭제", "제거"},
	Short:   "Remove Rss feed",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		if err := DelRssService(Ctx, serviceName); err != nil {
			log.Print(err)
			Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("failed to del")).Do()
			return
		}
		Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("success")).Do()
	},
}

var cmdRssList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"목록", "리스트"},
	Short:   "list Rss feed",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		//Ctx.UserID = `U80c288156ed296cfa61e8325df0e271c`
		msg, err := FlexContainerRssServices(Ctx, Ctx.UserID)
		if err != nil {
			log.Print(err)
			Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("failed to list")).Do()
			return
		}
		if _, err := Ctx.Linebot.ReplyMessage(
			Ctx.LineEvent.ReplyToken,
			linebot.NewFlexMessage(`alt`, msg),
		).Do(); err != nil {
			log.Print(err)
			return
		}
		//if _, err := Ctx.Linebot.PushMessage(
		//	Ctx.UserID,
		//	linebot.NewFlexMessage(`alt`, msg),
		//).Do(); err != nil {
		//	log.Print(err)
		//	return
		//}
	},
}

var cmdRssSubscribe = &cobra.Command{
	Use:     "subscribe [name]",
	Aliases: []string{"구독"},
	Short:   "Subscribe Rss feed",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		//Ctx.UserID = `U80c288156ed296cfa61e8325df0e271c`
		if err := SubscribeRssService(Ctx, serviceName, Ctx.UserID); err != nil {
			log.Print(err)
			Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("failed to subscribe")).Do()
			return
		}
		Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("success")).Do()
	},
}

var cmdRssUnsubscribe = &cobra.Command{
	Use:     "unsubscribe [name]",
	Aliases: []string{"구독해제"},
	Short:   "Unsubscribe Rss feed",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		//Ctx.UserID = `U80c288156ed296cfa61e8325df0e271c`
		if err := UnsubscribeRssService(Ctx, serviceName, Ctx.UserID); err != nil {
			log.Print(err)
			Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("failed to unsubscribe")).Do()
			return
		}
		Ctx.Linebot.ReplyMessage(Ctx.LineEvent.ReplyToken, linebot.NewTextMessage("success")).Do()
	},
}
