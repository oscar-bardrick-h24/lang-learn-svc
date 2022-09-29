package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/OJOMB/lang-learn-svc/internal/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_successCase(t *testing.T) {
	// create semi random user
	id, err := idTool.New()
	assert.NoError(t, err)

	var (
		email      = fmt.Sprintf("%s@example.com", id)
		password   = "password123"
		firstName  = id[:4]
		lastName   = id[4:]
		profilePic = "https://example.com/profilepic.jpg"
	)

	newUser := createUserReq{
		Email:      email,
		Password:   password,
		FirstName:  firstName,
		LastName:   lastName,
		ProfilePic: profilePic,
	}
	newUserBytes, err := json.Marshal(newUser)
	assert.NoError(t, err)

	tgtURL := fmt.Sprintf("%s/v1/users", getBaseUrl())
	req, err := http.NewRequest("POST", tgtURL, bytes.NewReader(newUserBytes))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respUser domain.User
	err = json.Unmarshal(respBytes, &respUser)
	assert.NoError(t, err)

	assert.Equal(t, newUser.Email, respUser.Email)
	assert.Equal(t, newUser.FirstName, respUser.FirstName)
	assert.Equal(t, newUser.LastName, respUser.LastName)
	assert.Equal(t, newUser.ProfilePic, respUser.ProfilePic)

	assert.True(t, idTool.IsValid(respUser.ID))
}
