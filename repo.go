package p

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type RssServiceRepoFirestore struct {
	Client *firestore.Client
}

var ErrClientNil = errors.New("firestore client is nil")

func (repo *RssServiceRepoFirestore) Get(sname string) (RssService, error) {
	var service RssService
	if repo.Client == nil {
		return service, ErrClientNil
	}

	ctx := context.Background()
	snap, err := repo.Client.Collection("RssService").Doc(sname).Get(ctx)
	if err != nil {
		return service, err
	}
	if err := snap.DataTo(&service); err != nil {
		return service, err
	}
	return service, nil
}

func (repo *RssServiceRepoFirestore) Add(service RssService) error {
	if repo.Client == nil {
		return ErrClientNil
	}
	//service.ID = String2sha256(service.Url)
	service.ID = service.Name
	service.UpdatedAt = time.Now()

	ctx := context.Background()
	_, err := repo.Client.Collection("RssService").Doc(service.ID).Set(ctx, service)
	return err
}

func (repo *RssServiceRepoFirestore) Del(service RssService) error {
	if repo.Client == nil {
		return ErrClientNil
	}
	//service.ID = String2sha256(service.Url)
	service.ID = service.Name

	ctx := context.Background()
	_, err := repo.Client.Collection("RssService").Doc(service.ID).Delete(ctx)
	return err
}

func (repo *RssServiceRepoFirestore) Foreach(handle func(service RssService)) error {
	if repo.Client == nil {
		return ErrClientNil
	}

	wg := sync.WaitGroup{}
	ctx := context.Background()
	iter := repo.Client.Collection("RssService").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return err
		}

		var service RssService
		if err := doc.DataTo(&service); err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			handle(service)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

type RssItemRepoFirestore struct {
	Client *firestore.Client
}

func (repo *RssItemRepoFirestore) Get(id string) (RssItem, error) {
	var item RssItem
	if repo.Client == nil {
		return item, ErrClientNil
	}

	ctx := context.Background()
	snap, err := repo.Client.Collection("RssItem").Doc(String2sha256(id)).Get(ctx)
	if err != nil {
		return item, err
	}
	if err := snap.DataTo(&item); err != nil {
		return item, err
	}
	return item, nil
}

func (repo *RssItemRepoFirestore) Add(item RssItem) error {
	if repo.Client == nil {
		return ErrClientNil
	}
	item.ID = String2sha256(item.Link)

	ctx := context.Background()
	_, err := repo.Client.Collection("RssItem").Doc(item.ID).Set(ctx, item)
	return err
}

func (repo *RssItemRepoFirestore) Del(item RssItem) error {
	if repo.Client == nil {
		return ErrClientNil
	}

	item.ID = String2sha256(item.Link)

	ctx := context.Background()
	_, err := repo.Client.Collection("RssItem").Doc(item.ID).Delete(ctx)
	return err
}

func (repo *RssItemRepoFirestore) IsNewItem(item RssItem) bool {
	_, err := repo.Get(item.Link)
	if err == nil {
		return false
	}

	stat, ok := status.FromError(err)
	if !ok || stat.Code() != codes.NotFound {
		return false
	}
	return true
}

func (repo *RssItemRepoFirestore) For(attr, op string, value interface{}, handle func(item RssItem) error) error {
	if repo.Client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	iter := repo.Client.Collection("RssItem").Where(attr, op, value).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return err
		}

		var item RssItem
		if err := doc.DataTo(&item); err != nil {
			return err
		}
		if err := handle(item); err != nil {
			return err
		}
	}

	return nil
}

type SubscriptionRepoFirestore struct {
	Client *firestore.Client
}

func (repo *SubscriptionRepoFirestore) Add(sub Subscription) error {
	if repo.Client == nil {
		return ErrClientNil
	}

	sub.ID = sub.ServiceName + `.` + sub.UserID

	ctx := context.Background()
	_, err := repo.Client.Collection("Subscription").Doc(sub.ID).Set(ctx, sub)
	return err
}

func (repo *SubscriptionRepoFirestore) Del(sub Subscription) error {
	if repo.Client == nil {
		return ErrClientNil
	}

	sub.ID = sub.ServiceName + `.` + sub.UserID

	ctx := context.Background()
	_, err := repo.Client.Collection("Subscription").Doc(sub.ID).Delete(ctx)
	return err
}

func (repo *SubscriptionRepoFirestore) Foreach(attr, op, value string, handle func(sub Subscription) error) error {
	if repo.Client == nil {
		return ErrClientNil
	}

	wg := sync.WaitGroup{}
	ctx := context.Background()
	iter := repo.Client.Collection("Subscription").Where(attr, op, value).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return err
		}

		var sub Subscription
		if err := doc.DataTo(&sub); err != nil {
			return err
		}
		//if err := handle(sub); err != nil {
		//	return err
		//}
		wg.Add(1)
		go func() {
			handle(sub)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (repo *SubscriptionRepoFirestore) ForSubscriber(sname string, handle func(sub Subscription) error) error {
	return repo.Foreach("service_name", "==", sname, handle)
}

func (repo *SubscriptionRepoFirestore) AllSubsByUser(uid string, handle func(sub Subscription) error) error {
	return repo.Foreach("user_id", "==", uid, handle)
}
