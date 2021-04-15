/*
token_builder.go
Implements the TokenBuilder to generate User Tokens

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
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenBuilderInterface defines the interface for token builders
type TokenBuilderInterface interface {
	CreateUserToken(user *User) (*UserTokenData, error)
}

// UserTokenData holds all information about a user token.
type UserTokenData struct {
	AccessToken  string
	RefreshToken string
	ATExpiresAt  time.Time
	RFExpiresAt  time.Time
}

// TokenBuilder implements TokenbuilderInterface
type TokenBuilder struct{}

// CreateUserToken builds a new UserTokenData instance for given user
func (t *TokenBuilder) CreateUserToken(user *User) (*UserTokenData, error) {
	atSecret := os.Getenv("AUTH_ACCESS_SECRET")
	rfSecret := os.Getenv("AUTH_REFRESH_SECRET")

	var err error

	// build TokenData
	td := &UserTokenData{
		ATExpiresAt: time.Now().Add(time.Minute * 15).UTC(),
		RFExpiresAt: time.Now().Add(time.Hour * 24 * 7).UTC(),
	}

	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID

	permissions, err := json.Marshal(user.Permissions)
	if err != nil {
		return nil, err
	}
	atClaims["permissions"] = string(permissions)

	atClaims["exp"] = td.ATExpiresAt
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(atSecret))
	if err != nil {
		return nil, err
	}

	// Create Refresh Token
	rfClaims := jwt.MapClaims{}
	rfClaims["authorized"] = true
	rfClaims["user_id"] = user.ID
	rfClaims["exp"] = td.RFExpiresAt
	rf := jwt.NewWithClaims(jwt.SigningMethodHS256, rfClaims)
	td.RefreshToken, err = rf.SignedString([]byte(rfSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}
