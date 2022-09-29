package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/OJOMB/lang-learn-svc/internal/pkg/domain"
	"github.com/OJOMB/lang-learn-svc/internal/pkg/uuidv4"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randLangCode() string {
	b := make([]rune, 2)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

const (
	localHost = "0.0.0.0"
	testPort  = 8080
)

var (
	idTool = uuidv4.NewGenerator()
)

type createUserReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
}

type createCourseReq struct {
	Title   string            `json:"title"`
	Lessons []createLessonReq `json:"lessons"`
}

type createLessonReq struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

type authReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createLanguageReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func getBaseUrl() string {
	return fmt.Sprintf("http://%s:%d", localHost, testPort)
}

func createNewUserAndAuthenticate(t *testing.T) (id string, token string) {
	// create semi random user
	id, err := idTool.New()
	assert.NoError(t, err)

	var (
		email      = fmt.Sprintf("%s@example.com", id)
		password   = "password"
		firstName  = id[:4]
		lastName   = id[4:8]
		profilePic = fmt.Sprintf("https://example.com/%s/profilepic.jpg", id)
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

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respUser domain.User
	err = json.Unmarshal(respBodyBytes, &respUser)
	assert.NoError(t, err)

	// authenticate user
	auth := authReq{
		Email:    email,
		Password: password,
	}
	authBytes, err := json.Marshal(auth)
	assert.NoError(t, err)

	tgtURL = fmt.Sprintf("%s/auth", getBaseUrl())
	req, err = http.NewRequest("POST", tgtURL, bytes.NewReader(authBytes))
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respAuth struct {
		Token string `json:"access_token"`
	}
	err = json.Unmarshal(respBodyBytes, &respAuth)
	assert.NoError(t, err)

	return respUser.ID, respAuth.Token
}
