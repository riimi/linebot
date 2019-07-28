package p

import (
	"github.com/mmcdole/gofeed"
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
	return repoRssService.Del(RssService{Name: name})
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
