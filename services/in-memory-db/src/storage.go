package main

import (
	"log"
	"time"
)

type StorageInterface interface {
	Initialize()
	Get(realmName string, key string) (bool, *Value)
	Set(realmName string, key string, value *Value)
	Delete(realmName string, key string) bool
	Keys(realmName string) []string
	Realms() []string
}

type Storage struct {
	Data map[string]map[string]*Value
}

func (s *Storage) Initialize() {
	s.Data = make(map[string]map[string]*Value)
}

func (s *Storage) GetRealm(realm string) (bool, map[string]*Value) {
	if _, ok := s.Data[realm]; ok {
		return true, s.Data[realm]
	}

	return false, make(map[string]*Value)
}

func (s *Storage) CreateRealm(realm string) map[string]*Value {
	if _, ok := s.Data[realm]; ok {
		return s.Data[realm]
	}

	s.Data[realm] = make(map[string]*Value)

	return s.Data[realm]
}

func (s *Storage) CleanEmptyRealm(realmName string) {
	if realm, ok := s.Data[realmName]; ok {
		if len(realm) == 0 {
			delete(s.Data, realmName)
		}
	}
}

func (s *Storage) Get(realmName string, key string) (bool, *Value) {
	ok, realm := s.GetRealm(realmName)
	if !ok {
		return false, nil
	}

	if val, ok := realm[key]; ok {
		return true, val
	}

	return false, nil
}

func (s *Storage) Set(realmName string, key string, value *Value) {
	ok, realm := s.GetRealm(realmName)
	if !ok {
		realm = s.CreateRealm(realmName)
	}

	realm[key] = value

	expiresIn := (int)(value.ExpiresAt.Sub(time.Now().UTC()).Seconds())

	expireFunc := func() {
		delete(realm, key)
		s.CleanEmptyRealm(realmName)
		log.Printf("Deleted key %v after %v seconds\n", key, expiresIn)
	}

	time.AfterFunc(time.Duration(expiresIn)*time.Second, expireFunc)
	log.Printf("Set key %v. It will Expire in %v seconds\n", key, expiresIn)
}

func (s *Storage) Delete(realmName string, key string) bool {
	ok, realm := s.GetRealm(realmName)
	if !ok {
		return false
	}

	if _, ok := realm[key]; ok {
		delete(realm, key)
		s.CleanEmptyRealm(realmName)
		return true
	}

	return false
}

func (s *Storage) Keys(realmName string) []string {
	ok, realm := s.GetRealm(realmName)
	if !ok {
		return make([]string, 0)
	}

	keys := make([]string, 0, len(realm))
	for k := range realm {
		keys = append(keys, k)
	}

	return keys
}

func (s *Storage) Realms() []string {
	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}

	return keys
}
