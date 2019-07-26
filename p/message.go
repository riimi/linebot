package p

import "github.com/line/line-bot-sdk-go/linebot"

func SourceID(source linebot.EventSource) string {
	switch source.Type {
	case linebot.EventSourceTypeUser:
		return source.UserID
	case linebot.EventSourceTypeGroup:
		return source.GroupID
	case linebot.EventSourceTypeRoom:
		return source.RoomID
	default:
		return ""
	}
}
