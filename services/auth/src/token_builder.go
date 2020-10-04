package main

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	guuid "github.com/google/uuid"
)

type TokenData struct {
	AccessToken      string
	RefreshToken     string
	AccessTokenUUID  string
	RefreshTokenUUID string
	ATExpiresAt      time.Time
	RFExpiresAt      time.Time
}

type TokenBuilder struct{}

func (t *TokenBuilder) CreateToken(user *User) (*TokenData, error) {
	atSecret := os.Getenv("AUTH_ACCESS_SECRET")
	rfSecret := os.Getenv("AUTH_REFRESH_SECRET")

	var err error

	// build TokenData
	td := &TokenData{
		ATExpiresAt:      time.Now().Add(time.Minute * 15).UTC(),
		AccessTokenUUID:  guuid.New().String(),
		RFExpiresAt:      time.Now().Add(time.Hour * 24 * 7).UTC(),
		RefreshTokenUUID: guuid.New().String(),
	}

	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_token_uuid"] = td.AccessTokenUUID
	atClaims["user_id"] = user.ID
	atClaims["permissions"] = user.Permissions
	atClaims["exp"] = td.ATExpiresAt
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(atSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}
