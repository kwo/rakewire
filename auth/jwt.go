package auth

import (
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	claimUserID         = "userid"
	claimName           = "name"
	claimRoles          = "roles"
	jwtExpiration       = time.Hour
	jwtSigningKeyLength = 64
	letters             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	jwtSigningMethod = crypto.SigningMethodHS256
	signingKey       []byte
	signingKeyLock   sync.Mutex
)

// GenerateToken creates a new JWT encoded with the username and roles of the given user, returns the token and expiration
func GenerateToken(user *User) ([]byte, time.Time, error) {
	expiration := time.Now().Add(jwtExpiration)
	claims := jws.Claims{}
	claims.SetExpiration(float64(expiration.Unix()))
	claims.Set(claimUserID, user.ID)
	claims.Set(claimName, user.Name)
	claims.Set(claimRoles, strings.Join(user.Roles, " "))
	token := jws.NewJWT(claims, jwtSigningMethod)
	tokenBytes, err := token.Serialize(getSigningKey())
	return tokenBytes, expiration, err
}

// RegenerateSigningKey generates the signing key
func RegenerateSigningKey() error {
	signingKeyLock.Lock()
	defer signingKeyLock.Unlock()
	signingKey = []byte(randomKey(jwtSigningKeyLength))
	return nil
}

func authenticateJWT(authHeader string, roles ...string) (*User, error) {

	tokenString := strings.TrimPrefix(authHeader, schemeJWT)
	token, errParse := jws.ParseJWT([]byte(tokenString))
	if errParse != nil {
		return nil, ErrBadHeader
	}

	if err := token.Validate(getSigningKey(), jwtSigningMethod); err != nil {
		return nil, ErrUnauthenticated
	}

	// test for missing expiration; already expired tokens are validated above in token.Validate
	if _, ok := token.Claims().Expiration(); !ok {
		return nil, ErrUnauthenticated
	}

	user := &User{}
	if claim := token.Claims().Get(claimUserID); claim != nil {
		if id, ok := claim.(string); ok {
			user.ID = id
		}
	}
	if claim := token.Claims().Get(claimName); claim != nil {
		if name, ok := claim.(string); ok {
			user.Name = name
		}
	}
	if claim := token.Claims().Get(claimRoles); claim != nil {
		if roles, ok := claim.(string); ok {
			user.Roles = strings.Fields(roles)
		}
	}

	if len(user.ID) > 0 {
		return user, nil
	}

	return nil, ErrUnauthenticated

}

func getSigningKey() []byte {
	signingKeyLock.Lock()
	defer signingKeyLock.Unlock()
	// generate new if empty
	if len(signingKey) == 0 {
		signingKey = []byte(randomKey(jwtSigningKeyLength))
	}
	return signingKey
}

func randomKey(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
