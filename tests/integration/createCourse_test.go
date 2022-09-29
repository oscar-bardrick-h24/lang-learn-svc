package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/OJOMB/lang-learn-svc/internal/pkg/domain"
	"github.com/stretchr/testify/assert"
)

// TestCreateCourse_successCase tests the happy path for creating a course
// obviously not the nicest code here but hopefully demonstrates the idea
func TestCreateCourse_successCase(t *testing.T) {
	var (
		courseTitle  = "courseTitleTest123"
		lesson1Title = "lesson1TitleTest123"
		lesson1Text  = "lesson1TextTest123"
		lesson2Title = "lesson2TitleTest123"
		lesson2Text  = "lesson2TextTest123"
		lang         = "GO"
	)

	// create user and authenticate to receive JWT token for use in subsequent requests
	userID, token := createNewUserAndAuthenticate(t)

	// create language
	langCreate := createLanguageReq{
		Code: randLangCode(),
		Name: "golang",
	}

	langCreateBytes, err := json.Marshal(langCreate)
	assert.NoError(t, err)

	tgtURL := fmt.Sprintf("%s/v1/languages", getBaseUrl())
	req, err := http.NewRequest("POST", tgtURL, bytes.NewReader(langCreateBytes))
	assert.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// create course
	newCourse := createCourseReq{
		Title: courseTitle,
		Lessons: []createLessonReq{
			{
				Title:    lesson1Title,
				Text:     lesson1Text,
				Language: lang,
			},
			{
				Title:    lesson2Title,
				Text:     lesson2Text,
				Language: lang,
			},
		},
	}

	newCourseBytes, err := json.Marshal(newCourse)
	assert.NoError(t, err)

	tgtURL = fmt.Sprintf("%s/v1/courses", getBaseUrl())
	req, err = http.NewRequest("POST", tgtURL, bytes.NewReader(newCourseBytes))
	assert.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	beforeCreation := time.Now()

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	afterCreation := time.Now()

	// time to check the created course is as expected
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respCourse domain.Course
	err = json.Unmarshal(respBytes, &respCourse)
	assert.NoError(t, err)

	assert.Equal(t, newCourse.Title, respCourse.Title)
	assert.Equal(t, userID, respCourse.CreatedBy)

	assert.True(t, idTool.IsValid(respCourse.ID))

	assert.Equal(t, len(newCourse.Lessons), len(respCourse.LessonIDs))
	for _, lessonID := range respCourse.LessonIDs {
		assert.True(t, idTool.IsValid(lessonID))
	}

	// now check the created lessons are as expected
	for _, lessonID := range respCourse.LessonIDs {
		tgtURL = fmt.Sprintf("%s/v1/lessons/%s", getBaseUrl(), lessonID)
		req, err = http.NewRequest("GET", tgtURL, nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBytes, err = ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respLesson domain.Lesson
		err = json.Unmarshal(respBytes, &respLesson)
		assert.NoError(t, err)

		assert.Equal(t, userID, respLesson.CreatedBy)
		assert.True(t, idTool.IsValid(respLesson.ID))
		assert.True(t, respLesson.CreatedAt.After(beforeCreation))
		assert.True(t, respLesson.CreatedAt.Before(afterCreation))
		assert.True(t, respLesson.UpdatedAt.After(beforeCreation))
		assert.True(t, respLesson.UpdatedAt.Before(afterCreation))

		if respLesson.Title == lesson1Title {
			assert.Equal(t, lesson1Text, respLesson.Text)
			assert.Equal(t, lang, respLesson.Language)
		} else if respLesson.Title == lesson2Title {
			assert.Equal(t, lesson2Text, respLesson.Text)
			assert.Equal(t, lang, respLesson.Language)
		} else {
			assert.Fail(t, "unexpected lesson title")
		}
	}
}
