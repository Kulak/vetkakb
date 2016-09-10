package core

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// User plays central role in the security model.
type User struct {
	UserID     int64
	Clearances uint8
	Created    time.Time
	Updated    time.Time
}

// OAuthUser contains data returned by OAuth provider.
// It points to vetka User.
type OAuthUser struct {
	UserFK            int64
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	ProvUserID        string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
	Created           time.Time
	Updated           time.Time
}

// Clearance represents access permission mask.
type Clearance struct {
	Mask uint8
	Name string
}

// Administrator access permission mask.
var Administrator = &Clearance{0x8, "Administrator"}

// Guest access permission mask.
var Guest = &Clearance{0x1, "Gust"}

// HasAccess returns true if clearance matches mask.
func (c *Clearance) HasAccess(clearances uint8) bool {
	return (c.Mask & clearances) > 0
}

// HasClearance returns true if user has specified clearance.
func (u *User) HasClearance(c *Clearance) bool {
	return c.HasAccess(u.Clearances)
}

func (u *User) dbInsert(tx *sql.Tx) (err error) {
	if u.UserID != 0 {
		return fmt.Errorf("Cannot insert User record, because UserID is not zero, but %v.",
			u.UserID)
	}
	query := `insert into user`
	var result sql.Result
	result, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("Cannot insert user record.  Error: %v", err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf("Failed to get affected User records after insert to DB. Error: %v", err)
	}
	u.UserID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to get inserted UserID. Error: %v", err)
	}
	return
}

func (ou *OAuthUser) dbInsert(tx *sql.Tx) (err error) {
	var result sql.Result
	var query string
	if ou.UserFK == 0 {
		return fmt.Errorf("Cannot insert oauthUser record without existing UserFK %v.", ou.UserFK)
	}
	query = `
insert into oauthUser
(userFK, provider, email, name, firstName, lastName, nickName,
description, provUserID, avatarURL, location, accessToken,
accessTokenSecret, refreshToken, expiresAt)
values
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
$11, $12, $13, $14, $15)
`
	result, err = tx.Exec(query,
		ou.UserFK, ou.Provider, ou.Email, ou.Name, ou.FirstName, ou.LastName, ou.NickName,
		ou.Description, ou.ProvUserID, ou.AvatarURL, ou.Location, ou.AccessToken,
		ou.AccessTokenSecret, ou.RefreshToken, ou.ExpiresAt)
	if err != nil {
		return fmt.Errorf("Failed to insert oauthUser record. Error: %v", err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf("Failed to get affected oauthUser records after insert to DB. Error: %v",
			err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to get EntryID for last insert operation. Error: %v", err)
	}
	if lastID != ou.UserFK {
		log.Printf("lastId != ou.UserFK: %v != %v", lastID, ou.UserFK)
	}
	return err
}
