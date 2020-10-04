package main

import "time"

type Message struct {
	ID              string                 `json:"string"` // TODO generate from timestamp (human read able and sortable) + topic + sender + uuid == Filename
	TopicIdentifier string                 `json:"topic"`
	Sender          string                 `json:"sender"`
	PayLoad         map[string]interface{} `json:"payload"`
	Deliveries      []*Delivery            `json:"delivery"`
	CreatedAt       time.Time              `json:"created-at"`
}

type Delivery struct {
	SubscriptionUrl string    `json:"subscription-url"`
	Delivered       bool      `json:"delivered"`
	Failed          bool      `json:"failed"`
	DeliveredAt     time.Time `json:"delivered-at"`
	FailedAt        time.Time `json:"failed-at"`
}
