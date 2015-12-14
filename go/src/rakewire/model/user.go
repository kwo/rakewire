package model

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

//User defines a system user
type User struct {
	ID           uint64
	Username     string
	PasswordHash string
	FeverHash    string
}

const (
	fUsername     = "Username"
	fPasswordHash = "PasswordHash"
	fFeverHash    = "FeverHash"
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
	return "User"
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

	if z.ID != 0 {
		result[fID] = strconv.FormatUint(z.ID, 10)
	}

	if z.Username != empty {
		result[fUsername] = z.Username
	}

	if z.PasswordHash != empty {
		result[fPasswordHash] = z.PasswordHash
	}

	if z.FeverHash != empty {
		result[fFeverHash] = z.FeverHash
	}

	return result

}

// Deserialize serializes an object to a list of key-values.
func (z *User) Deserialize(values map[string]string) error {

	for k, v := range values {
		switch k {
		case fID:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			z.ID = id
		case fUsername:
			z.Username = v
		case fPasswordHash:
			z.PasswordHash = v
		case fFeverHash:
			z.FeverHash = v
		}
	}

	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *User) IndexKeys() map[string][]string {
	result := make(map[string][]string)
	result[fUsername] = []string{strings.ToLower(z.Username)}
	result[fFeverHash] = []string{z.FeverHash}
	return result
}
