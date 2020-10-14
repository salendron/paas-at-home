package main

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	guuid "github.com/google/uuid"
)


type TokenBuilderInterface {
	
}

type UserTokenData struct {
	AccessToken      string
	RefreshToken     string
	ATExpiresAt      time.Time
	RFExpiresAt      time.Time
}

type ServiceTokenData struct {
	AccessToken string
}

type TokenBuilder struct{}

func (t *TokenBuilder) CreateUserToken(user *User) (*UserTokenData, error) {
	atSecret := os.Getenv("AUTH_ACCESS_SECRET")
	rfSecret := os.Getenv("AUTH_REFRESH_SECRET")

	var err error

	// build TokenData
	td := &UserTokenData{
		ATExpiresAt:      time.Now().Add(time.Minute * 15).UTC(),
		RFExpiresAt:      time.Now().Add(time.Hour * 24 * 7).UTC(),
	}

	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID
	atClaims["permissions"] = user.Permissions
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

func (t *TokenBuilder) CreateServiceToken(service *Service) (*ServiceTokenData,  error) {
	atSecret := os.Getenv("AUTH_SERVICE_SECRET")

	var err error

	td := ServiceTokenData{}

	// Create Token
	sClaims := jwt.MapClaims{}
	sClaims["authorized"] = true
	sClaims["service_id"] = service.ID
	s := jwt.NewWithClaims(jwt.SigningMethodHS256, sClaims)
	td.AccessToken, err = s.SignedString([]byte(atSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}