package domain

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (mr *mockRepo) CreateUser(ctx context.Context, user User) error {
	args := mr.Called(ctx, user)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) GetUser(ctx context.Context, userID string) (*User, error) {
	args := mr.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (mr *mockRepo) UpdateUser(ctx context.Context, user User) error {
	args := mr.Called(ctx, user)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := mr.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (mr *mockRepo) GetUserCourses(ctx context.Context, userID string) ([]UserCourse, error) {
	args := mr.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]UserCourse), args.Error(1)
}

func (mr *mockRepo) DeleteUser(ctx context.Context, userID string) error {
	args := mr.Called(ctx, userID)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) EnrollUser(ctx context.Context, userID, courseID string) error {
	args := mr.Called(ctx, userID, courseID)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) CreateLanguage(ctx context.Context, language Language) error {
	args := mr.Called(ctx, language)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) GetLanguages(ctx context.Context) ([]Language, error) {
	args := mr.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Language), args.Error(1)
}

func (mr *mockRepo) UpdateLanguage(ctx context.Context, lang Language) error {
	args := mr.Called(ctx, lang)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) DeleteLanguage(ctx context.Context, code string) error {
	args := mr.Called(ctx, code)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) DeleteLanguages(ctx context.Context) error {
	args := mr.Called(ctx)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) CreateCourse(ctx context.Context, course Course, lessons ...Lesson) error {
	args := mr.Called(ctx, course, lessons)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) GetCourse(ctx context.Context, courseID string) (*Course, error) {
	args := mr.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Course), args.Error(1)
}

func (mr *mockRepo) GetCourses(ctx context.Context) ([]Course, error) {
	args := mr.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Course), args.Error(1)
}

func (mr *mockRepo) GetCoursesByCreator(ctx context.Context, createdBy string) ([]Course, error) {
	args := mr.Called(ctx, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Course), args.Error(1)
}

func (mr *mockRepo) UpdateCourse(ctx context.Context, course Course) error {
	args := mr.Called(ctx, course)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) DeleteCourse(ctx context.Context, courseID string) error {
	args := mr.Called(ctx, courseID)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) AppendNewLessonToCourse(ctx context.Context, courseID string, lesson Lesson) error {
	args := mr.Called(ctx, courseID, lesson)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) CreateLesson(ctx context.Context, lesson Lesson) error {
	args := mr.Called(ctx, lesson)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) CreateLessons(ctx context.Context, lessons []Lesson) error {
	args := mr.Called(ctx, lessons)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) GetLesson(ctx context.Context, lessonID string) (*Lesson, error) {
	args := mr.Called(ctx, lessonID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Lesson), args.Error(1)
}

func (mr *mockRepo) GetLessons(ctx context.Context) ([]Lesson, error) {
	args := mr.Called(ctx)

	var rv []Lesson
	if args.Get(0) != nil {
		rv = args.Get(0).([]Lesson)
	}

	return rv, args.Error(1)
}

func (mr *mockRepo) DeleteLesson(ctx context.Context, lessonID string) error {
	args := mr.Called(ctx, lessonID)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (mr *mockRepo) UpdateLesson(ctx context.Context, lesson Lesson) error {
	args := mr.Called(ctx, lesson)
	if args.Error(0) == nil {
		return nil
	}

	return args.Error(0)
}

type mockIDTool struct {
	mock.Mock
}

func (mIDt *mockIDTool) New() (string, error) {
	args := mIDt.Called()

	return args.Get(0).(string), args.Error(1)
}

func (mIDt *mockIDTool) IsValid(ID string) bool {
	args := mIDt.Called(ID)
	return args.Get(0).(bool)
}

type mockPasswordTool struct {
	mock.Mock
}

func (mpt *mockPasswordTool) New(password string) (string, error) {
	args := mpt.Called(password)

	return args.Get(0).(string), args.Error(1)
}

func (mpt *mockPasswordTool) IsValid(password string) bool {
	args := mpt.Called(password)
	return args.Get(0).(bool)
}

func (mpt *mockPasswordTool) Check(hash, password string) error {
	args := mpt.Called(hash, password)
	return args.Error(0)
}

// Creates a silent logger instance that discards all output
func nullLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	return logger
}
