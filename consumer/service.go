package consumer

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

func AddRssService(repo *RssServiceRepoFirestore, name, url string) error {
	if repo == nil {
		return ErrClientNil
	}

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

func SubscribeService(repo *SubscriptionRepoFirestore, sname, uid string) error {
	if repo == nil {
		return ErrClientNil
	}

	sub, err := repo.Get(sname)
	if err != nil {
		return err
	}
	sub.AddSubscriber(uid)
	if err := repo.Add(sub); err != nil {
		return err
	}

	return nil
}

func UnsubscribeService(repo *SubscriptionRepoFirestore, sname, uid string) error {
	if repo == nil {
		return ErrClientNil
	}

	sub, err := repo.Get(sname)
	if err != nil {
		return err
	}
	sub.RemoveSubscriber(uid)
	if err := repo.Add(sub); err != nil {
		return err
	}

	return nil
}
