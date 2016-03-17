package model

//go:generate gokv $GOFILE

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

//User defines a system user
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username" kv:"+required,+groupby,Username:1:lower"`
	PasswordHash string `json:"passwordhash"`
	FeverHash    string `json:"feverhash" kv:"FeverHash:1"`
}

// NewUser creates a new user with the specified username
func NewUser(username string) *User {
	return &User{
		Username: username,
	}
}

func (z *User) setID(fn fnUniqueID) error {
	if z.ID == empty {
		if id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}

// SetPassword updates the password hashes.
func (z *User) SetPassword(password string) error {

	hash := md5.New()
	hash.Write([]byte(z.Username))
	hash.Write([]byte(":"))
	hash.Write([]byte(password))
	z.FeverHash = hex.EncodeToString(hash.Sum(nil))

	bhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	z.PasswordHash = hex.EncodeToString(bhash)

	return nil

}

// MatchPassword checks if the given password matches the user password.
func (z *User) MatchPassword(password string) bool {
	hashedPassword, err := hex.DecodeString(z.PasswordHash)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}
