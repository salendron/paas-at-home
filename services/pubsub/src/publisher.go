package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	guuid "github.com/google/uuid"
)

type PublisherInterface interface {
	Initialize(storage StorageInterface)
	Publish(messageType *MessageType, topic *Topic, sender string) (string, int, error)
}

type Publisher struct {
	Storage       StorageInterface
	PublishWorker func(message *Message, storage StorageInterface, isNew bool)
}

//Initialize initializes the Publisher by setting the active storage
func (p *Publisher) Initialize(storage StorageInterface) {
	p.Storage = storage
	p.PublishWorker = func(message *Message, storage StorageInterface, isNew bool) {
		log.Printf("Delivering message %v...", message.ID)
		hadErrors := false

		for i := 0; i < len(message.Deliveries); i++ {
			payload, err := json.Marshal(message.PayLoad)
			if err != nil {
				log.Printf("ERROR: Could not marshal payload of message %v\n", message.ID)
				hadErrors = true
				continue
			}

			resp, err := http.Post(message.Deliveries[i].SubscriptionUrl, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("ERROR: Could not deliver message %v to %v. Err: %v\n", message.ID, message.Deliveries[i].SubscriptionUrl, err)
				hadErrors = true

				message.Deliveries[i].Failed = true
				message.Deliveries[i].FailedAt = time.Now().UTC()
			} else {
				if resp.StatusCode != 200 {
					log.Printf("ERROR: Could not deliver message %v to %v. StatusCode: %v\n", message.ID, message.Deliveries[i].SubscriptionUrl, resp.StatusCode)
					hadErrors = true

					message.Deliveries[i].Failed = true
					message.Deliveries[i].FailedAt = time.Now().UTC()
				} else {
					log.Printf("SUCCESS: Message %v delivered to %v\n", message.ID, message.Deliveries[i].SubscriptionUrl)
					message.Deliveries[i].Delivered = true
					message.Deliveries[i].DeliveredAt = time.Now().UTC()
				}
			}

			err = storage.SaveMessage(message, isNew)
			if err != nil {
				log.Printf("ERROR: Could not save message %v\n", message.ID)
				hadErrors = true
			}
		}

		if isNew || (!isNew && !hadErrors) {
			err := storage.MoveMessage(message, isNew, !hadErrors)
			if err != nil {
				log.Printf("ERROR: Could not move message %v to correct state. Err: %v\n", message.ID, err)
			}
		}

		log.Printf("Finished delivering message %v. Had Errors: %v\n", message.ID, hadErrors)
	}

	// start worker to watch for failed messages
	p.runRetryFailedDeliveriesWorker()
}

func (p *Publisher) Publish(messageType *MessageType, topic *Topic, sender string) (string, int, error) {
	if len(topic.Subscriptions) == 0 {
		log.Printf("Topic %v has no subscribers. Message from sender %v will not be processed\n", topic.Identifier, sender)
		return "", 0, nil
	}

	message := &Message{
		ID: strings.Join(
			[]string{time.Now().UTC().Format(time.RFC3339Nano), topic.Identifier, sender, guuid.New().String()},
			"-",
		),
		TopicIdentifier: topic.Identifier,
		Sender:          sender,
		PayLoad:         messageType.Message,
		Deliveries:      make([]*Delivery, len(topic.Subscriptions)),
		CreatedAt:       time.Now().UTC(),
	}

	// build deliveries
	for i := 0; i < len(topic.Subscriptions); i++ {
		message.Deliveries = append(message.Deliveries,
			&Delivery{
				SubscriptionUrl: topic.Subscriptions[i].URL,
				Delivered:       false,
				Failed:          false,
			},
		)
	}

	// save message
	err := p.Storage.SaveMessage(message, true)
	if err != nil {
		return message.ID, len(topic.Subscriptions), err
	}

	p.runPublishWorker(message, true)

	return message.ID, len(topic.Subscriptions), nil
}

func (p *Publisher) runRetryFailedDeliveriesWorker() {
	worker := func(storage StorageInterface) {
		log.Println("Starting to watch for failed messages to retry.")
		for {

			messages, err := storage.ListMessages(false)
			log.Printf("Will retry %v failed messages", len(messages))

			if err != nil {
				log.Printf("ERROR: Could not load messages to retry. Err: %v\n", err)
			} else {
				for i := 0; i < len(messages); i++ {
					p.PublishWorker(messages[i], storage, false)
				}
			}

			// wait
			time.Sleep(5 * time.Minute)
		}
	}

	go worker(p.Storage)
}

func (p *Publisher) runPublishWorker(message *Message, isNew bool) {
	go p.PublishWorker(message, p.Storage, isNew)
}
