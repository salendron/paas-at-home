package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Subscription struct {
	URL                    string    `json:"url"`
	CreatedAt              time.Time `json:"created-at"`
	LastSuccessfulDelivery time.Time `json:"last-successful-delivery"`
	FailedDeliveries       int       `json:"failed-deliveries"`
	LastFailedDelivery     time.Time `json:"last-failed-delivery"`
}

type Topic struct {
	Identifier     string          `json:"identifier"`
	AllowedSenders []string        `json:"allowed-senders"`
	Subscriptions  []*Subscription `json:"subscriptions"`
}

func (t *Topic) Validate() error {
	t.Identifier = strings.ToUpper(strings.TrimSpace(t.Identifier))

	if len(t.Identifier) == 0 {
		return errors.New("Topic Identifier can't be empty")
	}

	if t.AllowedSenders == nil {
		t.AllowedSenders = make([]string, 0)
	}

	return nil
}

func (t *Topic) IsSenderAllowed(sender string) bool {
	for i := 0; i < len(t.AllowedSenders); i++ {
		if t.AllowedSenders[i] == sender {
			return true
		}
	}

	return false
}

// Subscribe adds a new Subscription to a topic.
// returns true and an error on error
// returns false and an error if sender is not allowed
// returns true and no error if successful
func (t *Topic) Subscribe(url string, storage StorageInterface) (bool, error) {
	for i := 0; i < len(t.Subscriptions); i++ {
		if t.Subscriptions[i].URL == url {
			return true, nil // URL is already subscribed
		}
	}

	// create new subscription
	sub := &Subscription{}
	sub.URL = url
	sub.CreatedAt = time.Now().UTC()
	t.Subscriptions = append(t.Subscriptions, sub)

	err := storage.SaveTopic(t)
	if err != nil {
		return true, err
	}

	return true, nil
}

// Unsubscribe removes a Subscription from a topic.
func (t *Topic) Unsubscribe(url string, storage StorageInterface) error {
	newSubscriptions := make([]*Subscription, 0)
	for i := 0; i < len(t.Subscriptions); i++ {
		if t.Subscriptions[i].URL != url {
			newSubscriptions = append(newSubscriptions, t.Subscriptions[i])
		}
	}
	t.Subscriptions = newSubscriptions

	return storage.SaveTopic(t)
}

// TEST
func Test() {
	subscriptions := make(chan string)

	TestChannel(subscriptions)

	WriteToChan("1", subscriptions)
	time.Sleep(2 * time.Second)
	WriteToChan("2", subscriptions)
	time.Sleep(3 * time.Second)
	WriteToChan("3", subscriptions)
	WriteToChan("4", subscriptions)
}

func WriteToChan(s string, c chan string) {
	c <- s
}

func TestChannel(c chan string) {
	go func(c chan string) {
		for {
			fmt.Println(<-c)
		}
	}(c)
}
