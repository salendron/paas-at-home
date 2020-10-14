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

var usersDirectory = "users"
var applicationsDirectory = "applications"

//StorageInterface defines the interface for the data storage.
type StorageInterface interface {
	GetUserByCredentials(username string, passowrd string) (*User, bool, error)
	GetUser(ID string) (*User, error)
	ListUSers([]*User, error)
	SaveUser(user *User) error
	DeleteUser(user *User) error
	GetUserByCredentials(username string, password string) (*User, bool, error)
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

// writeFile writes a file
func (s *Storage) writeFile(filePath string, file []byte) error {
	s.MutexLock.Lock()
	defer s.MutexLock.Unlock()

	return ioutil.WriteFile(filePath, file, 0755)
}

// GetUser loads a user. If it does not exists it returns nil as user
func (s *Storage) GetUser(ID string) (*User, error) {
	usersPath, err := s.getDirectoryPath(usersDirectory)
	if err != nil {
		return nil, err
	}

	userPath := filepath.Join(usersPath, fmt.Sprintf("%v.json", ID))
	user := &User{}
	if s.fileExists(userPath) {
		jsonFile, err := os.Open(userPath)
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, user)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return user, nil
}

// ListUsers lists all users
func (s *Storage) ListUsers() ([]*User, error) {
	users := make([]*User, 0)

	usersPath, err := s.getDirectoryPath(usersDirectory)
	if err != nil {
		return nil, err
	}

	identifiers := make([]string, 0)
	filepath.Walk(usersPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".json", f.Name())
			if err == nil && r {
				identifiers = append(identifiers, strings.TrimRight(f.Name(), ".json"))
			}
		}
		return nil
	})

	for i := 0; i < len(identifiers); i++ {
		user, err := s.GetUser(identifiers[i])
		if err != nil {
			return users, nil
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) SaveUser(user *User) error {
	usersPath, err := s.getDirectoryPath(usersDirectory)
	if err != nil {
		return err
	}

	userPath := filepath.Join(usersPath, fmt.Sprintf("%v.json", user.ID))

	file, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = s.writeFile(userPath, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(user *User) error {
	usersPath, err := s.getDirectoryPath(usersDirectory)
	if err != nil {
		return err
	}

	userPath := filepath.Join(usersPath, fmt.Sprintf("%v.json", user.ID))

	return os.Remove(userPath)
}

func (s *Storage) GetUserByCredentials(username string, password string) (*User, bool, error) {
	user, err := s.GetUser(username)
	if err != nil {
		return nil, false, err
	}

	if user.Password != password {
		return nil, false, nil
	}

	return user, true, nil
}
