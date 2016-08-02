package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Roles        []string `json:"roles"`
	PasswordHash string   `json:"passwordhash"`
	FeverHash    string   `json:"feverhash"`
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

// AddRole adds the given role to the user.
func (z *User) AddRole(role string) {
	if !z.HasRole(role) {
		z.Roles = append(z.Roles, role)
	}
}

// HasRole tests if the user has been assigned the given role
func (z *User) HasRole(role string) bool {
	result := false
	if len(z.Roles) > 0 {
		for _, value := range z.Roles {
			if value == role {
				return true
			}
		}
	}
	return result
}

// HasAllRoles tests if the user has all the given roles
func (z *User) HasAllRoles(roles ...string) bool {
	for _, role := range roles {
		if !z.HasRole(role) {
			return false
		}
	}
	return true
}

// RemoveRole removes a role from the user.
func (z *User) RemoveRole(role string) {
	for i, value := range z.Roles {
		if value == role {
			z.Roles = append(z.Roles[:i], z.Roles[i+1:]...)
		}
	}
}

// RoleString formats all roles as a comma-separated string
func (z *User) RoleString() string {
	var result string
	for i, role := range z.Roles {
		if i > 0 {
			result = result + ", "
		}
		result = result + role
	}
	return result
}

// SetRoles replaces roles with the roles contained withing the given comma-separated string.
func (z *User) SetRoles(rolestr string) {
	roles := strings.Split(rolestr, ",")
	z.Roles = []string{}
	for _, role := range roles {
		z.AddRole(strings.TrimSpace(role))
	}
}

func (z *User) clear() {
	z.ID = empty
	z.Username = empty
	z.PasswordHash = empty
	z.FeverHash = empty
	z.Roles = []string{}
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

func (z *User) hasIncrementingID() bool {
	return true
}

func (z *User) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexUserUsername] = []string{strings.ToLower(z.Username)}
	result[indexUserFeverhash] = []string{z.FeverHash}
	return result
}

func (z *User) setID(tx Transaction) error {
	id, err := tx.NextID(entityUser)
	if err != nil {
		return err
	}
	z.ID = keyEncodeUint(id)
	return nil
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

func (z *Users) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Users) encode() ([]byte, error) {
	return json.Marshal(z)
}
