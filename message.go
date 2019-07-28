package p

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
)

func SourceID(source *linebot.EventSource) string {
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

func BubbleContainerRssService(service RssService, nowsub bool) *linebot.BubbleContainer {
	OneFlex := 1
	FiveFlex := 5

	ret := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		//		Hero: &linebot.ImageComponent{
		//			Type:        linebot.FlexComponentTypeImage,
		//			URL:         service.Url,
		//			Size:        linebot.FlexImageSizeTypeFull,
		//			AspectRatio: linebot.FlexImageAspectRatioType20to13,
		//			AspectMode:  linebot.FlexImageAspectModeTypeCover,
		//		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   service.Name,
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "Title",
									Flex:  &OneFlex,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#aaaaaa",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  service.Title,
									Flex:  &FiveFlex,
									Size:  linebot.FlexTextSizeTypeSm,
									Wrap:  true,
									Color: "#666666",
								},
							},
							Spacing: linebot.FlexComponentSpacingTypeSm,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "Desc",
									Flex:  &OneFlex,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#aaaaaa",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  service.Description,
									Flex:  &FiveFlex,
									Size:  linebot.FlexTextSizeTypeSm,
									Wrap:  true,
									Color: "#666666",
								},
							},
							Spacing: linebot.FlexComponentSpacingTypeSm,
						},
					},
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Margin:  linebot.FlexComponentMarginTypeLg,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.PostbackAction{
						Label:       "Subscribe",
						Data:        fmt.Sprintf(`!rss subscribe %s`, service.Name),
						DisplayText: fmt.Sprintf(`!rss subscribe %s`, service.Name),
					},
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypePrimary,
				},
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.URIAction{
						Label: "Website",
						URI:   service.Url,
					},
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypeSecondary,
				},
				&linebot.SpacerComponent{
					Type: linebot.FlexComponentTypeSpacer,
					Size: linebot.FlexSpacerSizeTypeSm,
				},
			},
			Spacing: linebot.FlexComponentSpacingTypeSm,
		},
	}
	if nowsub {
		ret.Footer.Contents[0] = &linebot.ButtonComponent{
			Type: linebot.FlexComponentTypeButton,
			Action: &linebot.PostbackAction{
				Label:       "Unsubscribe",
				Data:        fmt.Sprintf(`!rss unsubscribe %s`, service.Name),
				DisplayText: fmt.Sprintf(`!rss unsubscribe %s`, service.Name),
			},
			Height: linebot.FlexButtonHeightTypeSm,
			Style:  linebot.FlexButtonStyleTypeSecondary,
		}
	}

	return ret
}

func FlexContainerRssServices(ctx *Context, uid string) (linebot.FlexContainer, error) {
	services := make([]RssService, 0)
	subsMap := make(map[string]bool)
	repoRssService := &RssServiceRepoFirestore{Client: ctx.Firestore}
	if err := repoRssService.Foreach(func(serv RssService) error {
		services = append(services, serv)
		subsMap[serv.Name] = false
		return nil
	}); err != nil {
		return nil, err
	}

	repoSubs := &SubscriptionRepoFirestore{Client: ctx.Firestore}
	if err := repoSubs.AllSubsByUser(uid, func(sub Subscription) error {
		subsMap[sub.ServiceName] = true
		return nil
	}); err != nil {
		return nil, err
	}

	contents := make([]*linebot.BubbleContainer, 0, len(subsMap))
	for _, serv := range services {
		contents = append(contents, BubbleContainerRssService(serv, subsMap[serv.Name]))
	}
	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: contents,
	}, nil
}
