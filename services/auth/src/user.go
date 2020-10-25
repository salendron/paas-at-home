package main

// User contains all information about a user to login
type User struct {
	ID          string // Equals Username - has to be unique anyway
	Password    string
	Permissions []Permission
}
