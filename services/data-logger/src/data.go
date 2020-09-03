/*
data.go
Implements a single data instance, which is saved to a collection storage.
It also takes care of initializing the Data instance with a UUID and a created-at
date.

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
	"io"
	"strings"
	"time"

	guuid "github.com/google/uuid"
)

// DataFileContent holds all Data items in a single data file
type DataFileContent struct {
	Items []*Data `json:"items"`
}

// GetItemsInRange returns all items of a data file between startDate and endDate.
func (df DataFileContent) GetItemsInRange(startDate time.Time, endDate time.Time) []*Data {
	items := make([]*Data, 0)

	for _, item := range df.Items {
		if item.CreatedAt.After(startDate) || item.CreatedAt.Equal(startDate) || item.CreatedAt.Before(endDate) || item.CreatedAt.Equal(endDate) {
			items = append(items, item)
		}
	}

	return items
}

// Data wrapps the actual logged data (payload) in a normalized sstructure containing,
// an UUID, a CreatedAt date and also the original data that should be saved (Payload).
type Data struct {
	UUID      string                 `json:"uuid"`
	CreatedAt time.Time              `json:"created-at"`
	Payload   map[string]interface{} `json:"payload"`
}

// Initialize sets the UUID and CreatedAt date.
func (d *Data) Initialize() {
	d.UUID = strings.Join(
		[]string{time.Now().UTC().Format("2006-01-02"), guuid.New().String()},
		"-",
	)
	d.CreatedAt = time.Now().UTC()
}

// PayloadFromRequestJson parses the given json data of the request's io.ReadCloser
func PayloadFromRequestJson(rc io.ReadCloser) (map[string]interface{}, error) {
	payload := make(map[string]interface{})

	err := json.NewDecoder(rc).Decode(&payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
