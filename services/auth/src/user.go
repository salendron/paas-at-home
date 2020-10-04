package main

type User struct {
	ID          string // Equals Username - has to be unique anyway
	Password    string
	Permissions []*Permission
}
