package model

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

//User defines a system user
type User struct {
	ID           string `kv:"Username:2,FeverHash:2"`
	Username     string `kv:"Username:1"`
	PasswordHash string
	FeverHash    string `kv:"FeverHash:1"`
}

// NewUser creates a new user with the specified username
func NewUser(username string) *User {
	return &User{
		ID:       getUUID(),
		Username: username,
	}
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
