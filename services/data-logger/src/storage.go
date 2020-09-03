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
	"sync"
	"time"
)

//StorageInterface defines the interface for the data storage.
type StorageInterface interface {
	Initialize(dataRootDirectory string)
	ReadData(collectionName string, startDate time.Time, endDate time.Time) ([]*Data, error)
	WriteData(collectionName string, payload map[string]interface{}) (*Data, error)
	ListCollections() ([]string, error)
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

// getCollectionPath get's the actual path of a collection.
// If the collection does not exist it will create a new directory for the
// colelction first.
func (s *Storage) getCollectionPath(collectionName string) (string, error) {
	path := filepath.Join(s.DataRootDirectory, collectionName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return path, err
		}
	}

	return path, nil
}

// getCurrentDatafilePath gets the path of the current data file
func (s *Storage) getCurrentDatafilePath(collectionName string) (string, error) {
	collectionPath, err := s.getCollectionPath(collectionName)
	if err != nil {
		return "", err
	}

	currentDataFileName := fmt.Sprintf("%v.%v", time.Now().UTC().Format("2006-01-02"), "json")

	return filepath.Join(collectionPath, currentDataFileName), nil
}

// getDataFilePathsInRange get's all path of the data files of a collection
// used in a given time range.
func (s *Storage) getDataFilePathsInRange(collectionName string, startDate time.Time, endDate time.Time) ([]string, error) {
	dataFilePaths := make([]string, 0)

	// truncate time part to 00:00:00
	currentDate := startDate.Truncate(24 * time.Hour)
	limitDate := endDate.Truncate(24 * time.Hour)

	// get collaction path
	collectionPath, err := s.getCollectionPath(collectionName)
	if err != nil {
		return dataFilePaths, err
	}

	for currentDate.Before(limitDate) || currentDate.Equal(limitDate) {
		dataFilePaths = append(
			dataFilePaths,
			filepath.Join(
				collectionPath,
				fmt.Sprintf("%v.%v", currentDate.Format("2006-01-02"), "json"),
			),
		)

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return dataFilePaths, nil
}

// fileExists checks if a file exists
func (s *Storage) fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// readDataFile reads the content of a data file to a *DataFileContent
func (s *Storage) readDataFile(dataFilePath string) (*DataFileContent, error) {
	data := &DataFileContent{}

	if s.fileExists(dataFilePath) {
		jsonFile, err := os.Open(dataFilePath)
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, data)
		if err != nil {
			return nil, err
		}
	} else {
		data.Items = make([]*Data, 0)
	}

	return data, nil
}

// ReadData loads all data items of a collection in a given time range
func (s *Storage) ReadData(collectionName string, startDate time.Time, endDate time.Time) ([]*Data, error) {
	result := make([]*Data, 0)

	filePaths, err := s.getDataFilePathsInRange(collectionName, startDate, endDate)
	if err != nil {
		return result, err
	}

	for _, df := range filePaths {
		if s.fileExists(df) {
			dfContent, err := s.readDataFile(df)
			if err != nil {
				return result, err
			}

			result = append(result, dfContent.GetItemsInRange(startDate, endDate)...)
		}
	}

	return result, nil
}

// writeDataFile writes a data file
func (s *Storage) writeDataFile(dataFilePath string, file []byte) error {
	s.MutexLock.Lock()
	defer s.MutexLock.Unlock()

	return ioutil.WriteFile(dataFilePath, file, 0755)
}

// WriteData stores a new data item to a collection and returns it wrapped in a
// *Data structure
func (s *Storage) WriteData(collectionName string, payload map[string]interface{}) (*Data, error) {
	data := &Data{
		Payload: payload,
	}
	data.Initialize()

	dataFilePath, err := s.getCurrentDatafilePath(collectionName)
	if err != nil {
		return data, err
	}

	dataFile, err := s.readDataFile(dataFilePath)

	dataFile.Items = append(dataFile.Items, data)

	file, err := json.Marshal(dataFile)
	if err != nil {
		return data, err
	}

	err = s.writeDataFile(dataFilePath, file)
	if err != nil {
		return data, err
	}

	return data, nil
}

// ListCollections return all available collections
func (s *Storage) ListCollections() ([]string, error) {
	collections := make([]string, 0)

	dirContent, err := ioutil.ReadDir(s.DataRootDirectory)
	if err != nil {
		return collections, err
	}

	for _, item := range dirContent {
		if item.IsDir() {
			collections = append(collections, item.Name())
		}
	}

	return collections, nil
}
