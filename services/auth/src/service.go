package main

// Service contains all information about a service to login
type Service struct {
	ID      string
	AuthKey string // actually something like the password of this service
}
