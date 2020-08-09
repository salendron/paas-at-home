package main

import (
	"log"
	"time"
)

type StorageInterface interface {
	Initialize()
	Get() (bool, string, error)
	Set(key string, value string)
	Delete() error
	Keys() []string
}

type Storage struct {
	Values map[string]string
}

func (s *Storage) Initialize() {
	s.Values = make(map[string]string)
}

func (s *Storage) Set(key string, value string, expiresIn int) {
	s.Values[key] = value

	expireFunc := func() {
		delete(s.Values, key)
		log.Printf("Deleted key %v after %v seconds\n", key, expiresIn)
	}

	time.AfterFunc(time.Duration(expiresIn)*time.Second, expireFunc)
	log.Printf("Set key %v. It will Expire in %v seconds\n", key, expiresIn)
}
