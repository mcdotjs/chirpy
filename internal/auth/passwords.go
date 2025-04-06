package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// type MyCustomClaims struct {
// 	jwt.RegisteredClaims
// }
 type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}
// func MakeJWT3(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
// 	mySigningKey := []byte(tokenSecret)
//
// 	parsedID := userID.String()
// 	claims := MyCustomClaims{
// 		jwt.RegisteredClaims{
// 			Issuer:    "chirpy",
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			NotBefore: jwt.NewNumericDate(time.Now()),
// 			Subject:   parsedID,
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	ss, err := token.SignedString(mySigningKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	return ss, nil
// }

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}
// func ValidadteJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(tokenSecret), nil
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
// 		return uuid.MustParse(claims.RegisteredClaims.Subject), nil
// 	} else {
// 		fmt.Errorf("unknown claims type, cannot proceed")
// 	}
// 	return uuid.UUID{}, nil
// }
