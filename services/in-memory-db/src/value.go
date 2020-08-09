package main

import (
	"encoding/json"
	"io"
	"time"
)

type Value struct {
	Value     string
	ExpiresAt time.Time
}

func (v *Value) ToValueMessageType() ValueMessageType {
	return ValueMessageType{
		Value:     v.Value,
		ExpiresIn: (int)(v.ExpiresAt.Sub(time.Now().UTC()).Seconds()),
	}
}

func ValueFromValueMessageType(body io.ReadCloser) (*Value, error) {
	msg := ValueMessageType{}

	err := json.NewDecoder(body).Decode(&msg)
	if err != nil {
		return nil, err
	}

	value := &Value{
		Value:     msg.Value,
		ExpiresAt: time.Now().UTC().Add(time.Duration(msg.ExpiresIn) * time.Second),
	}

	return value, nil
}
