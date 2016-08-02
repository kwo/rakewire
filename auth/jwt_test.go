package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestJWT(t *testing.T) {

	user1 := &User{
		ID:    "0000000001",
		Name:  "jake",
		Roles: []string{"admin", "operator"},
	}

	tokenString, _, errGen := GenerateToken(user1)
	if errGen != nil {
		t.Errorf("Error generating token: %s", errGen.Error())
	}
	if len(tokenString) == 0 {
		t.Error("Nil token")
	}
	t.Logf("generated token string: %s", tokenString)

	user2, errAuth := authenticateJWT(schemeJWT + string(tokenString))
	if errAuth != nil {
		t.Errorf("Error authenticating token: %s", errAuth.Error())
	}
	if user2 == nil {
		t.Error("Nil auth user")
	}

	if user2.Name != user1.Name {
		t.Errorf("Bad username: %s, expected: %s", user2.Name, user1.Name)
	}

	if strings.Join(user2.Roles, ",") != strings.Join(user1.Roles, ",") {
		t.Errorf("Bad roles: %v, expected: %v", user2.Roles, user1.Roles)
	}

}

func TestNoExp(t *testing.T) {

	claims := jwt.MapClaims{}
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	claims[claimName] = "jake"
	tokenString, errSign := token.SignedString(getSigningKey())
	if errSign != nil {
		t.Fatalf("Failed to sign token: %s", errSign.Error())
	}

	user, errAuth := authenticateJWT(schemeJWT + string(tokenString))
	if errAuth != ErrUnauthenticated {
		t.Errorf("Bad error: %v, expected %v", errAuth, ErrUnauthenticated)
	}
	if user != nil {
		t.Errorf("Bad user: %v, expected nil", user)
	}

}

func TestExpired(t *testing.T) {

	claims := jwt.MapClaims{}
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	claims[claimName] = "jake"
	claims[claimExpiration] = time.Now().Add(time.Second * -1).Unix()
	tokenString, errSign := token.SignedString(getSigningKey())
	if errSign != nil {
		t.Fatalf("Failed to sign token: %s", errSign.Error())
	}

	user, errAuth := authenticateJWT(schemeJWT + string(tokenString))
	if errAuth != ErrUnauthenticated {
		t.Errorf("Bad error: %v, expected %v", errAuth, ErrUnauthenticated)
	}
	if user != nil {
		t.Errorf("Bad user: %v, expected nil", user)
	}

}

func TestWrongSignature(t *testing.T) {

	claims := jwt.MapClaims{}
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	claims[claimName] = "jake"
	claims[claimExpiration] = time.Now().Add(time.Second * -1).Unix()
	tokenString, errSign := token.SignedString(getSigningKey())
	if errSign != nil {
		t.Fatalf("Failed to sign token: %s", errSign.Error())
	}

	RegenerateSigningKey()

	user, errAuth := authenticateJWT(schemeJWT + string(tokenString))
	if errAuth != ErrUnauthenticated {
		t.Errorf("Bad error: %v, expected %v", errAuth, ErrUnauthenticated)
	}
	if user != nil {
		t.Errorf("Bad user: %v, expected nil", user)
	}

}

func TestNoUser(t *testing.T) {

	claims := jwt.MapClaims{}
	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	claims[claimExpiration] = time.Now().Add(time.Second * -1).Unix()
	tokenString, errSign := token.SignedString(getSigningKey())
	if errSign != nil {
		t.Fatalf("Failed to sign token: %s", errSign.Error())
	}

	user, errAuth := authenticateJWT(schemeJWT + string(tokenString))
	if errAuth != ErrUnauthenticated {
		t.Errorf("Bad error: %v, expected %v", errAuth, ErrUnauthenticated)
	}
	if user != nil {
		t.Errorf("Bad user: %v, expected nil", user)
	}

}
