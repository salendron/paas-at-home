/*
storage.go
Defines the storage interface of this application and implements the
actual database object storage.

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
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var usersDirectory = "users"
var servicesDirectory = "services"

//StorageInterface defines the interface for the data storage.
type StorageInterface interface {
	GetUserByCredentials(username string, passowrd string) (*User, error)
	GetUserByID(ID int) (*User, error)
}

// Storage implements StorageInterface
type Storage struct {
	DB *gorm.DB
}

// Connect connects the storage to the database
func (s *Storage) Connect(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//models
	db.AutoMigrate(&User{})

	s.DB = db

	return err
}

// GetUserByID loads a user. If it does not exists it returns nil as user
func (s *Storage) GetUserByID(ID int) (*User, error) {
	user := &User{}
	result := s.DB.First(user, ID)

	err := result.Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, nil
}

// GetUserByCredentials loads a User using given credentials
func (s *Storage) GetUserByCredentials(username string, password string) (*User, error) {
	user := &User{}
	result := s.DB.Where("username = ?", username).Where("password = ?", password).First(user)
	err := result.Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
