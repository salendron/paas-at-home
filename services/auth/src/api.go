/*
api.go
Defines the RESTful API interface of this application and implements all api
methods.

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

/*
//APIInterface defines the interface of the RESTful API
type APIInterface interface {
	Initialize(storage StorageInterface, publisher PublisherInterface)
	CreateTopic(w http.ResponseWriter, r *http.Request)
	UpdateTopic(w http.ResponseWriter, r *http.Request)
	GetTopic(w http.ResponseWriter, r *http.Request)
	DeleteTopic(w http.ResponseWriter, r *http.Request)
	ListTopics(w http.ResponseWriter, r *http.Request)
	Subscribe(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	Publish(w http.ResponseWriter, r *http.Request)
}

//API implements APIInterface
type API struct {
	Storage   StorageInterface
	Publisher PublisherInterface
}

//Initialize initializes the API by setting the active storage and publisher
func (a *API) Initialize(storage StorageInterface, publisher PublisherInterface) {
	a.Storage = storage
	a.Publisher = publisher
}

// parseRequestPayload parses the given json data of the request's io.ReadCloser
func parseRequestPayload(rc io.ReadCloser, dst interface{}) error {
	err := json.NewDecoder(rc).Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

//API handler to create topics
func (a *API) CreateTopic(w http.ResponseWriter, r *http.Request) {
	topic := &Topic{}
	err := parseRequestPayload(r.Body, topic)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	err = topic.Validate()
	if err != nil {
		RaiseError(w, fmt.Sprintf("Invalid request body. %v", err.Error()), http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	existingTopic, err := a.Storage.GetTopic(topic.Identifier)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if existingTopic != nil {
		RaiseError(w, "Topic already exists", http.StatusBadRequest, ErrorCodeTopicExists)
		return
	}

	// initialize subscriptions
	topic.Subscriptions = make([]*Subscription, 0)

	err = a.Storage.SaveTopic(topic)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

//API handler to update topics
func (a *API) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, ok := vars["id"]
	if !ok {
		RaiseError(w, "ID is missing", http.StatusBadRequest, ErrorCodeIDIsMissing)
		return
	}

	topic := &Topic{}
	err := parseRequestPayload(r.Body, topic)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	err = topic.Validate()
	if err != nil {
		RaiseError(w, fmt.Sprintf("Invalid request body. %v", err.Error()), http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	ID = strings.ToUpper(strings.TrimSpace(ID))

	if ID != topic.Identifier {
		RaiseError(w, "ID in request URL does not match topic identifier", http.StatusBadRequest, ErrorCodeIDMismatch)
		return
	}

	existingTopic, err := a.Storage.GetTopic(topic.Identifier)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if existingTopic == nil {
		RaiseError(w, "Topic does not exists", http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	// subscriptions are read only in this call - use subscribe or unsubscribe to edit this
	topic.Subscriptions = existingTopic.Subscriptions

	err = a.Storage.SaveTopic(topic)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topic)
}

func (a *API) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, ok := vars["id"]
	if !ok {
		RaiseError(w, "ID is missing", http.StatusBadRequest, ErrorCodeIDIsMissing)
		return
	}

	ID = strings.ToUpper(strings.TrimSpace(ID))

	topic, err := a.Storage.GetTopic(ID)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if topic == nil {
		RaiseError(w, "Topic does not exists", http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	err = a.Storage.DeleteTopic(topic)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, ok := vars["id"]
	if !ok {
		RaiseError(w, "ID is missing", http.StatusBadRequest, ErrorCodeIDIsMissing)
		return
	}

	ID = strings.ToUpper(strings.TrimSpace(ID))

	topic, err := storage.GetTopic(ID)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusBadRequest, ErrorCodeInternal)
		return
	}

	if topic == nil {
		RaiseError(w, fmt.Sprintf("No topic found with ID: %v", ID), http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topic)
}

func (a *API) ListTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := storage.ListTopics()
	if err != nil {
		RaiseError(w, err.Error(), http.StatusBadRequest, ErrorCodeInternal)
		return
	}

	msg := TopicListMessageType{
		Topics: topics,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

//API handler to subscribe to a topic
func (a *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, ok := vars["topic"]
	if !ok {
		RaiseError(w, "Topic ID is missing", http.StatusBadRequest, ErrorCodeIDIsMissing)
		return
	}

	subscriptionUrl := &SubscriptionUrlType{}
	err := parseRequestPayload(r.Body, subscriptionUrl)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	subscriptionUrl.URL = strings.TrimSpace(subscriptionUrl.URL)

	if len(subscriptionUrl.URL) == 0 { // FIXME We should validate that this is a valid URL
		RaiseError(w, "Subscription Url can't be empty", http.StatusBadRequest, ErrorCodeInvalidSubscriptionURL)
		return
	}

	existingTopic, err := a.Storage.GetTopic(topicID)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if existingTopic == nil {
		RaiseError(w, "Topic does not exists", http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	// try to subscribe the URL
	senderAllowed, err := existingTopic.Subscribe(subscriptionUrl.URL, a.Storage)
	if !senderAllowed {
		RaiseError(w, "Sender is not allowed for this topic", http.StatusNotFound, ErrorCodeSenderNotAllowed)
		return
	}

	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

//API handler to unsubscribe from a task
func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, ok := vars["topic"]
	if !ok {
		RaiseError(w, "Topic ID is missing", http.StatusBadRequest, ErrorCodeIDIsMissing)
		return
	}

	subscriptionUrl := &SubscriptionUrlType{}
	err := parseRequestPayload(r.Body, subscriptionUrl)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	subscriptionUrl.URL = strings.TrimSpace(subscriptionUrl.URL)

	if len(subscriptionUrl.URL) == 0 { // FIXME We should validate that this is a valid URL
		RaiseError(w, "Subscription Url can't be empty", http.StatusBadRequest, ErrorCodeInvalidSubscriptionURL)
		return
	}

	existingTopic, err := a.Storage.GetTopic(topicID)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if existingTopic == nil {
		RaiseError(w, "Topic does not exists", http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	// Unsubscribe URL from topic
	err = existingTopic.Unsubscribe(subscriptionUrl.URL, a.Storage)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

//API handler to publish messages
func (a *API) Publish(w http.ResponseWriter, r *http.Request) {
	message := &MessageType{}
	err := parseRequestPayload(r.Body, message)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	message.TopicIdentifier = strings.ToUpper(strings.TrimSpace(message.TopicIdentifier))

	if len(message.TopicIdentifier) == 0 {
		RaiseError(w, "Topic can't be empty", http.StatusBadRequest, ErrorCodeInvalidTopic)
		return
	}

	topic, err := a.Storage.GetTopic(message.TopicIdentifier)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if topic == nil {
		RaiseError(w, "Topic does not exists", http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	messageID, subscriberCount, err := a.Publisher.Publish(message, topic, "sender")

	confirmation := MessageReceivedConfirmationType{
		ID:              messageID,
		TopicIdentifier: message.TopicIdentifier,
		SubscriberCount: subscriberCount,
		CreatedAt:       time.Now().UTC(),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(confirmation)
}
*/
