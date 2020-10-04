/*
storage.go
Defines the storage interface of this application and implements the
actual json object storage based on data files.

###################################################################################

MIT License

Copyright (c) 2020 Bruno Hautzenberger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var topicsDirectory = "topics"
var newMessagesDirectory = "new_messages"
var failedMessagesDirectory = "failed_messages"
var doneMessagesDirectory = "done_messages"

//StorageInterface defines the interface for the data storage.
type StorageInterface interface {
	Initialize(dataRootDirectory string)
	GetTopic(identifier string) (*Topic, error)
	ListTopics() ([]*Topic, error)
	SaveTopic(topic *Topic) error
	SaveMessage(message *Message, isNew bool) error
	DeleteTopic(topic *Topic) error
	MoveMessage(message *Message, wasNew bool, isDone bool) error
	ListMessages(isNew bool) ([]*Message, error)
	GetMessage(ID string, isNew bool) (*Message, error)
}

//Implements StorageInterface
type Storage struct {
	DataRootDirectory string
	MutexLock         sync.Mutex
}

// Initialize sets the data root directory
func (s *Storage) Initialize(dataRootDirectory string) {
	s.DataRootDirectory = dataRootDirectory
}

// fileExists checks if a file exists
func (s *Storage) fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// getDirectoryPath is used to get a directory inside the data root
// directory. If it does not exist, it will be created.
func (s *Storage) getDirectoryPath(name string) (string, error) {
	path := filepath.Join(s.DataRootDirectory, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return path, err
		}
	}

	return path, nil
}

// GetTopic loads a topic. If it does not exists it returns nil as topic
func (s *Storage) GetTopic(identifier string) (*Topic, error) {
	topicsPath, err := s.getDirectoryPath(topicsDirectory)
	if err != nil {
		return nil, err
	}

	topicPath := filepath.Join(topicsPath, fmt.Sprintf("%v.json", identifier))
	topic := &Topic{}
	if s.fileExists(topicPath) {
		jsonFile, err := os.Open(topicPath)
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, topic)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return topic, nil
}

// ListTopics lists all available topics
func (s *Storage) ListTopics() ([]*Topic, error) {
	topics := make([]*Topic, 0)

	topicsPath, err := s.getDirectoryPath(topicsDirectory)
	if err != nil {
		return nil, err
	}

	identifiers := make([]string, 0)
	filepath.Walk(topicsPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".json", f.Name())
			if err == nil && r {
				identifiers = append(identifiers, strings.TrimRight(f.Name(), ".json"))
			}
		}
		return nil
	})

	for i := 0; i < len(identifiers); i++ {
		topic, err := s.GetTopic(identifiers[i])
		if err != nil {
			return topics, nil
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

// writeFile writes a file
func (s *Storage) writeFile(filePath string, file []byte) error {
	s.MutexLock.Lock()
	defer s.MutexLock.Unlock()

	return ioutil.WriteFile(filePath, file, 0755)
}

// MoveMessage moves a message either to done or failed based on
// status of deliveries.
// wasNew indicates if this Ticket is in new message directory or if it is already in failed (failed already before)
// isDone indicates if this message was delivered to all subscribers or not
func (s *Storage) MoveMessage(message *Message, wasNew bool, isDone bool) error {
	// get src path
	srcSubPath := newMessagesDirectory
	if !wasNew {
		srcSubPath = failedMessagesDirectory
	}

	srcPath, err := s.getDirectoryPath(srcSubPath)
	if err != nil {
		return err
	}

	srcPath = filepath.Join(srcPath, fmt.Sprintf("%v.json", message.ID))

	// get dest path
	destSubPath := doneMessagesDirectory
	if !isDone {
		destSubPath = failedMessagesDirectory
	}

	destPath, err := s.getDirectoryPath(destSubPath)
	if err != nil {
		return err
	}

	destPath = filepath.Join(destPath, fmt.Sprintf("%v.json", message.ID))

	// do move
	return os.Rename(srcPath, destPath)
}

func (s *Storage) SaveTopic(topic *Topic) error {
	topicsPath, err := s.getDirectoryPath(topicsDirectory)
	if err != nil {
		return err
	}

	topicPath := filepath.Join(topicsPath, fmt.Sprintf("%v.json", topic.Identifier))

	file, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	err = s.writeFile(topicPath, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteTopic(topic *Topic) error {
	topicsPath, err := s.getDirectoryPath(topicsDirectory)
	if err != nil {
		return err
	}

	topicPath := filepath.Join(topicsPath, fmt.Sprintf("%v.json", topic.Identifier))

	return os.Remove(topicPath)
}

//
func (s *Storage) SaveMessage(message *Message, isNew bool) error {
	messagesSubDir := newMessagesDirectory
	if !isNew {
		messagesSubDir = failedMessagesDirectory
	}

	newMessagesPath, err := s.getDirectoryPath(messagesSubDir)
	if err != nil {
		return err
	}

	messagePath := filepath.Join(newMessagesPath, fmt.Sprintf("%v.json", message.ID))

	file, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = s.writeFile(messagePath, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListMessages(isNew bool) ([]*Message, error) {
	messages := make([]*Message, 0)

	messagesSubDir := newMessagesDirectory
	if !isNew {
		messagesSubDir = failedMessagesDirectory
	}

	messagesPath, err := s.getDirectoryPath(messagesSubDir)
	if err != nil {
		return nil, err
	}

	identifiers := make([]string, 0)
	filepath.Walk(messagesPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".json", f.Name())
			if err == nil && r {
				identifiers = append(identifiers, strings.TrimRight(f.Name(), ".json"))
			}
		}
		return nil
	})

	for i := 0; i < len(identifiers); i++ {
		message, err := s.GetMessage(identifiers[i], isNew)
		if err != nil {
			return messages, nil
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (s *Storage) GetMessage(ID string, isNew bool) (*Message, error) {
	messagesSubDir := newMessagesDirectory
	if !isNew {
		messagesSubDir = failedMessagesDirectory
	}

	messagesPath, err := s.getDirectoryPath(messagesSubDir)
	if err != nil {
		return nil, err
	}

	messagePath := filepath.Join(messagesPath, fmt.Sprintf("%v.json", ID))

	message := &Message{}
	if s.fileExists(messagePath) {
		jsonFile, err := os.Open(messagePath)
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, message)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return message, nil
}
