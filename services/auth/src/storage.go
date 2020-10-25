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
	"log"
	"os"
	"path/filepath"
)

var usersDirectory = "users"
var servicesDirectory = "services"

//StorageInterface defines the interface for the data storage.
type StorageInterface interface {
	Initialize(dataRootDirectory string)
	GetUserByCredentials(username string, passowrd string) (*User, bool, error)
	GetUser(ID string) (*User, error)
	GetServiceByCredentials(ID string, key string) (*Service, bool, error)
	GetService(ID string) (*Service, error)
}

// Storage implements StorageInterface
type Storage struct {
	DataRootDirectory string
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

// GetUser loads a user. If it does not exists it returns nil as user
func (s *Storage) GetUser(ID string) (*User, error) {
	if ID == "su" {
		user := &User{
			ID: "su",
			Permissions: []Permission{
				Permission{Key: "ROOT"},
			},
		}
		return user, nil
	}

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

// GetUserByCredentials loads a User using given credentials
func (s *Storage) GetUserByCredentials(username string, password string) (*User, bool, error) {
	log.Println(fmt.Sprintf("%v %v %v", username, password, os.Getenv("SU_PWD")))
	if username == "su" && password == os.Getenv("SU_PWD") {
		user := &User{
			ID: "su",
			Permissions: []Permission{
				Permission{Key: "ROOT"},
			},
		}
		return user, true, nil
	}
	user, err := s.GetUser(username)
	if err != nil {
		return nil, false, err
	}

	if user == nil {
		return nil, false, nil
	}

	if user.Password != password {
		return nil, false, nil
	}

	return user, true, nil
}

// GetService loads a service. If it does not exist it returns nil as service
func (s *Storage) GetService(ID string) (*Service, error) {
	storagePath, err := s.getDirectoryPath(servicesDirectory)
	if err != nil {
		return nil, err
	}

	servicePath := filepath.Join(storagePath, fmt.Sprintf("%v.json", ID))
	service := &Service{}
	if s.fileExists(servicePath) {
		jsonFile, err := os.Open(servicePath)
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, service)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return service, nil
}

// GetServiceByCredentials loads a Service using given credentials.
func (s *Storage) GetServiceByCredentials(ID string, key string) (*Service, bool, error) {
	service, err := s.GetService(ID)
	if err != nil {
		return nil, false, err
	}

	if service.AuthKey != key {
		return nil, false, nil
	}

	return service, true, nil
}
