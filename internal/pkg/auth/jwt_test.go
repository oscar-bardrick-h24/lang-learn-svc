package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTToolGenerateSignedTokenString_successPath(t *testing.T) {
	mIDT := &mockIDTool{}
	mIDT.On("New").Return("token_id", nil)

	jt := NewJWTTool("supersecretkey", time.Hour, "issuer", mIDT)

	tokenSigned, err := jt.GenerateTokenString("user_id")
	assert.NoError(t, err)
	assert.Regexp(t, `^(?:[\w-]*\.){2}[\w-]*$`, tokenSigned)

	tok, err := jwt.ParseWithClaims(tokenSigned, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jt.secretKey), nil
	})
	assert.NoError(t, err)

	claims, ok := tok.Claims.(*JWTClaims)
	assert.True(t, ok)
	assert.Equal(t, "issuer", claims.Issuer)
	assert.Equal(t, "user_id", claims.Subject)
}

func TestJWTToolGetSubject_successPath(t *testing.T) {
	mIDT := &mockIDTool{}
	mIDT.On("New").Return("token_id", nil)

	jt := NewJWTTool("supersecretkey", time.Hour, "issuer", mIDT)

	// JWT NumericDate doesn't deal in nanoseconds so we truncate to nearest second
	issuedAt := time.Now().Truncate(time.Second)
	expiresAt := issuedAt.Add(2 * time.Hour)
	userID := "user_id"
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(jt.secretKey)
	assert.NoError(t, err)

	subj, err := jt.GetSubject(signedStr)
	assert.NoError(t, err)
	assert.Equal(t, "user_id", subj)
}

func TestJWTToolGetSubject_expiredToken_failurePath(t *testing.T) {
	mIDT := &mockIDTool{}
	mIDT.On("New").Return("token_id", nil)

	jt := NewJWTTool("supersecretkey", time.Hour, "issuer", mIDT)

	// JWT NumericDate doesn't deal in nanoseconds so we truncate to nearest second
	issuedAt := time.Time{}
	expiresAt := issuedAt.Add(time.Hour)
	userID := "user_id"
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(jt.secretKey)
	assert.NoError(t, err)

	subj, err := jt.GetSubject(signedStr)
	assert.Equal(t, subj, "")
	assert.Error(t, err)

	assert.True(t, strings.HasPrefix(err.Error(), "failed to parse token: token is expired by"))
}
