package main

type ValueMessageType struct {
	Value     string `json:"value"`
	ExpiresIn int    `json:"expires-in"`
}

type KeyListMessageType struct {
	Keys []string `json:"keys"`
}

type RealmListMessageType struct {
	Realms []string `json:"realms"`
}

type ErrorMessageType struct {
	Error interface{} `json:"error"`
}
