package main

/*

type ArtistListMessageType struct {
	Data []*Artist `json:"data"`
}*/

type ValueMessageType struct {
	Value     string `json:"value"`
	ExpiresIn int    `json:"expires-in"`
}

type ErrorMessageType struct {
	Error interface{} `json:"error"`
}
