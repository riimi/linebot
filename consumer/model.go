package consumer

import (
	"crypto/sha256"
	"encoding/hex"
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
}

type RssItem struct {
	ID          string    `firestore:"id"`
	ServiceName string    `firestore:"service_name"`
	Title       string    `firestore:"title"`
	Link        string    `firestore:"link"`
	Published   string    `firestore:"published"`
	CreatedAt   time.Time `json:"-" firestore:"created_at"`
}

type Subscription struct {
	ID          string    `firestore:"id"`
	UserId      []string  `firestore:"user_id"`
	ServiceName string    `firestore:"service_name"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

func (sub *Subscription) AddSubscriber(uid string) {
	sub.UserId = append(sub.UserId, uid)
}

func (sub *Subscription) RemoveSubscriber(uid string) {
	deleted := -1
	for index, id := range sub.UserId {
		if id == uid {
			deleted = index
			break
		}
	}
	if deleted >= 0 {
		sub.UserId = append(sub.UserId[:deleted], sub.UserId[deleted+1:]...)
	}
}

func String2sha256(in string) string {
	sha := sha256.Sum256([]byte(in))
	return hex.EncodeToString(sha[:])
}
