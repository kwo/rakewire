package model

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//go:generate gokv $GOFILE

//User defines a system user
type User struct {
	ID           uint64
	Username     string `kv:"Username:1"`
	PasswordHash string
	FeverHash    string `kv:"FeverHash:1"`
}

// index names
const (
	UserEntity         = "User"
	UserIndexUsername  = "Username"
	UserIndexFeverHash = "FeverHash"
)

const (
	uID           = "ID"
	uUsername     = "Username"
	uPasswordHash = "PasswordHash"
	uFeverHash    = "FeverHash"
)

// NewUser creates a new user with the specified username
func NewUser(username string) *User {
	return &User{
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

// GetName return the name of the entity.
func (z *User) GetName() string {
	return UserEntity
}

// GetID return the primary key of the object.
func (z *User) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *User) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *User) Clear() {
	z.ID = 0
	z.Username = empty
	z.PasswordHash = empty
	z.FeverHash = empty
}

// Serialize serializes an object to a list of key-values.
func (z *User) Serialize() map[string]string {
	result := make(map[string]string)
	setUint(z.ID, uID, result)
	setString(z.Username, uUsername, result)
	setString(z.PasswordHash, uPasswordHash, result)
	setString(z.FeverHash, uFeverHash, result)
	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *User) Deserialize(values map[string]string) error {
	var errors []error
	z.ID = getUint(uID, values, errors)
	z.Username = getString(uUsername, values, errors)
	z.PasswordHash = getString(uPasswordHash, values, errors)
	z.FeverHash = getString(uFeverHash, values, errors)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *User) IndexKeys() map[string][]string {
	result := make(map[string][]string)
	result[UserIndexUsername] = []string{strings.ToLower(z.Username)}
	result[UserIndexFeverHash] = []string{z.FeverHash}
	return result
}
