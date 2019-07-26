package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type RssService struct {
	ID          string    `firestore:"id"`
	Name        string    `firestore:"name,omitempty"`
	Url         string    `firestore:"url,omitempty"`
	Title       string    `firestore:"title,omitempty"`
	Description string    `firestore:"description,omitempty"`
	Link        string    `firestore:"link,omitempty"`
	Language    string    `firestore:"language,omitempty"`
	NumFollwer  int       `firestore:"num_follower,omitempty"`
	CreatedAt   time.Time `firestore:"created_at,omitempty"`
	UpdatedAt   time.Time `firestore:"updated_at,omitempty"`
	UserID      []string  `firestore:"user_id"`
}

func (serv *RssService) AddSubscriber(uid string) {
	serv.UserID = append(serv.UserID, uid)
}

func (serv *RssService) RemoveSubscriber(uid string) {
	deleted := -1
	for index, id := range serv.UserID {
		if id == uid {
			deleted = index
			break
		}
	}
	if deleted >= 0 {
		serv.UserID = append(serv.UserID[:deleted], serv.UserID[deleted+1:]...)
	}
}

type RssItem struct {
	ID          string    `firestore:"id"`
	ServiceName string    `firestore:"service_name"`
	Title       string    `firestore:"title"`
	Link        string    `firestore:"link"`
	Published   string    `firestore:"published"`
	CreatedAt   time.Time `json:"-" firestore:"created_at"`
}

func (item *RssItem) TextMessage() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", item.ServiceName, item.Title, item.Published, item.Link)
}

func String2sha256(in string) string {
	sha := sha256.Sum256([]byte(in))
	return hex.EncodeToString(sha[:])
}
