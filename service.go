package p

import (
	"github.com/mmcdole/gofeed"
	"log"
	"os/exec"
	"time"
)

func GetRssFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	out, err := exec.Command("curl", "-c", "cookie", "-XGET", "-L", url).Output()
	if err != nil {
		return nil, err
	}
	feed, err := fp.ParseString(string(out))
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func AddRssService(ctx *Context, name, url string) error {
	if ctx == nil {
		return ErrClientNil
	}
	repo := RssServiceRepoFirestore{Client: ctx.Firestore}

	feed, err := GetRssFeed(url)
	if err != nil {
		return err
	}

	service := RssService{
		Name:        name,
		Url:         url,
		Title:       feed.Title,
		Description: feed.Description,
		Link:        feed.Link,
		Language:    feed.Language,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := repo.Add(service); err != nil {
		return err
	}

	return nil
}

func DelRssService(ctx *Context, name string) error {
	if ctx == nil {
		return ErrClientNil
	}

	repoRssService := &RssServiceRepoFirestore{Client: ctx.Firestore}
	if err := repoRssService.Del(RssService{Name: name}); err != nil {
		return err
	}
	tobedel := make([]RssItem, 0)
	repoRssItem := &RssItemRepoFirestore{Client: ctx.Firestore}
	if err := repoRssItem.For("service_name", "==", name, func(item RssItem) error {
		log.Print(item)
		tobedel = append(tobedel, item)
		return nil
	}); err != nil {
		return err
	}

	for _, item := range tobedel {
		if err := repoRssItem.Del(item); err != nil {
			return err
		}
	}
	return nil
}

func SubscribeRssService(ctx *Context, sname, uid string) error {
	if ctx == nil {
		return ErrClientNil
	}

	repoRssService := &RssServiceRepoFirestore{Client: ctx.Firestore}
	repoSub := SubscriptionRepoFirestore{Client: ctx.Firestore}

	_, err := repoRssService.Get(sname)
	if err != nil {
		return err
	}

	newSub := Subscription{
		ServiceName: sname,
		UserID:      uid,
		CreateAt:    time.Now(),
	}

	if err := repoSub.Add(newSub); err != nil {
		return err
	}

	return nil
}

func UnsubscribeRssService(ctx *Context, sname, uid string) error {
	if ctx == nil {
		return ErrClientNil
	}
	repoSub := SubscriptionRepoFirestore{Client: ctx.Firestore}

	return repoSub.Del(Subscription{
		ServiceName: sname,
		UserID:      uid,
	})
}

func CleanOldRssItem(ctx *Context) error {
	tobedel := make([]RssItem, 0)
	repo := &RssItemRepoFirestore{Client: ctx.Firestore}
	if err := repo.For("created_at", ">", time.Now().Add(time.Duration(-7)*time.Hour), func(item RssItem) error {
		log.Print(item)
		tobedel = append(tobedel, item)
		return nil
	}); err != nil {
		return err
	}

	for _, item := range tobedel {
		if err := repo.Del(item); err != nil {
			return err
		}
	}

	return nil
}
