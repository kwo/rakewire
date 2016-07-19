package auth

import (
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	claimUserID         = "userid"
	claimName           = "name"
	claimRoles          = "roles"
	claimExpiration     = "exp"
	jwtExpiration       = time.Hour
	jwtSigningKeyLength = 64
	letters             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
	signingKey       []byte
	signingKeyLock   sync.Mutex
)

// GenerateToken creates a new JWT encoded with the username and roles of the given user, returns the token and expiration
func GenerateToken(user *User) (string, time.Time, error) {
	expiration := time.Now().Add(jwtExpiration)
	claims := jwt.MapClaims{}
	claims[claimExpiration] = expiration.Unix()
	claims[claimUserID] = user.ID
	claims[claimName] = user.Name
	claims[claimRoles] = strings.Join(user.Roles, " ")
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	tokenString, err := token.SignedString(getSigningKey())
	return tokenString, expiration, err
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
	token, errParse := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthenticated
		}
		return getSigningKey(), nil
	})
	if errParse != nil {
		return nil, ErrUnauthenticated
	}

	user := &User{}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claim := claims[claimUserID]; claim != nil {
			if id, ok := claim.(string); ok {
				user.ID = id
			}
		}
		if claim := claims[claimName]; claim != nil {
			if name, ok := claim.(string); ok {
				user.Name = name
			}
		}
		if claim := claims[claimRoles]; claim != nil {
			if roles, ok := claim.(string); ok {
				user.Roles = strings.Fields(roles)
			}
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
