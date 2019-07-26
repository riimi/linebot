package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (repo *RssServiceRepoFirestore) Foreach(handle func(service RssService) error) error {
	if repo.Client == nil {
		return ErrClientNil
	}

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
		if err := handle(service); err != nil {
			return err
		}
	}

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
