package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	entityUser         = "User"
	indexUserUsername  = "Username"
	indexUserFeverhash = "Feverhash"
)

var (
	indexesUser = []string{
		indexUserFeverhash, indexUserUsername,
	}
)

// User defines a system user
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordhash"`
	FeverHash    string `json:"feverhash"`
}

// GetID returns the unique ID for the object
func (z *User) GetID() string {
	return z.ID
}

// MatchPassword checks if the given password matches the user password.
func (z *User) MatchPassword(password string) bool {
	hashedPassword, err := hex.DecodeString(z.PasswordHash)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
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

func (z *User) clear() {
	z.ID = empty
	z.Username = empty
	z.PasswordHash = empty
	z.FeverHash = empty
}

func (z *User) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *User) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *User) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexUserUsername] = []string{strings.ToLower(z.Username)}
	result[indexUserFeverhash] = []string{z.FeverHash}
	return result
}

func (z *User) setID(tx Transaction) error {
	seq := C.getSequences(tx)
	seq.User = seq.User + 1
	z.ID = keyEncodeUint(seq.User)
	return C.putSequences(tx, seq)
}

// Users is a collection of User objects
type Users []*User

// ByID maps items to their ID
func (z Users) ByID() map[string]*User {
	result := make(map[string]*User)
	for _, user := range z {
		result[user.ID] = user
	}
	return result
}

func (z Users) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z Users) encode() ([]byte, error) {
	return json.Marshal(z)
}
