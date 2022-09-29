package app

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/OJOMB/lang-learn-svc/internal/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (ms *mockService) CreateUser(ctx context.Context, email, pword, fname, lname, profPic string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, email, pword, fname, lname, profPic)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) GetUser(ctx context.Context, userID string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, userID)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) GetUserByUserName(ctx context.Context, userName string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, userName)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, email)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) GetUserCourses(ctx context.Context, userID string) ([]domain.UserCourse, *domain.Error) {
	args := ms.Called(ctx, userID)

	var courses []domain.UserCourse
	if args.Get(0) == nil {
		courses = nil
	} else {
		courses = args.Get(0).([]domain.UserCourse)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return courses, err
}

func (ms *mockService) EnrollUser(ctx context.Context, userID, courseID string) *domain.Error {
	args := ms.Called(ctx, userID, courseID)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) DeleteUser(ctx context.Context, userID string) *domain.Error {
	args := ms.Called(ctx, userID)

	var err *domain.Error
	if args.Get(0) == nil {
		err = nil
	} else {
		err = args.Get(0).(*domain.Error)
	}

	return err
}

func (ms *mockService) PatchUser(ctx context.Context, userID string, patchJSON []byte) *domain.Error {
	args := ms.Called(ctx, userID, patchJSON)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) ValidateUserCredentials(ctx context.Context, userName, email, pword string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, userName, email, pword)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) GetAuthenticatedUser(ctx context.Context, email, pword string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, email, pword)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) SetUserProfilePic(ctx context.Context, userID, profPic string) *domain.Error {
	args := ms.Called(ctx, userID, profPic)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)

}

func (ms *mockService) SetUserPassword(ctx context.Context, userID, password string) *domain.Error {
	args := ms.Called(ctx, userID, password)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)

}

func (ms *mockService) CreateLanguage(ctx context.Context, code, name string) *domain.Error {
	args := ms.Called(ctx, code, name)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)

}

func (ms *mockService) GetLanguages(ctx context.Context) ([]domain.Language, *domain.Error) {
	args := ms.Called(ctx)

	var languages []domain.Language
	if args.Get(0) != nil {
		languages = args.Get(0).([]domain.Language)
	}

	var err *domain.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*domain.Error)
	}

	return languages, err
}

func (ms *mockService) UpdateLanguage(ctx context.Context, code, name string) *domain.Error {
	args := ms.Called(ctx, code, name)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) DeleteLanguage(ctx context.Context, code string) *domain.Error {
	args := ms.Called(ctx, code)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) DeleteLanguages(ctx context.Context) *domain.Error {
	args := ms.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) CreateCourse(ctx context.Context, title string, lessons []domain.Lesson) (*domain.Course, *domain.Error) {
	args := ms.Called(ctx, title)

	var course *domain.Course
	if args.Get(0) == nil {
		course = nil
	} else {
		course = args.Get(0).(*domain.Course)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return course, err
}

func (ms *mockService) GetLesson(ctx context.Context, lessonID string) (*domain.Lesson, *domain.Error) {
	args := ms.Called(ctx, lessonID)

	var lesson *domain.Lesson
	if args.Get(0) == nil {
		lesson = nil
	} else {
		lesson = args.Get(0).(*domain.Lesson)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return lesson, err
}

func (ms *mockService) GetLessons(ctx context.Context) ([]domain.Lesson, *domain.Error) {
	args := ms.Called(ctx)

	var lessons []domain.Lesson
	if args.Get(0) != nil {
		lessons = args.Get(0).([]domain.Lesson)
	}

	var err *domain.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*domain.Error)
	}

	return lessons, err
}

func (ms *mockService) DeleteLesson(ctx context.Context, lessonID string) *domain.Error {
	args := ms.Called(ctx, lessonID)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) UpdateLesson(ctx context.Context, lessonID, title, text, language string) *domain.Error {
	args := ms.Called(ctx, lessonID, title, text, language)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) GetCourse(ctx context.Context, courseID string) (*domain.Course, *domain.Error) {
	args := ms.Called(ctx, courseID)

	var course *domain.Course
	if args.Get(0) == nil {
		course = nil
	} else {
		course = args.Get(0).(*domain.Course)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return course, err
}

func (ms *mockService) GetCourses(ctx context.Context, createdBy string) ([]domain.Course, *domain.Error) {
	args := ms.Called(ctx, createdBy)

	var courses []domain.Course
	if args.Get(0) != nil {
		courses = args.Get(0).([]domain.Course)
	}

	var err *domain.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*domain.Error)
	}

	return courses, err
}

func (ms *mockService) DeleteCourse(ctx context.Context, courseID string) *domain.Error {
	args := ms.Called(ctx, courseID)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

func (ms *mockService) AppendNewLessonToCourse(ctx context.Context, courseID, lessonTitle, lessonText, lessonLanguage string) (*domain.Lesson, *domain.Error) {
	args := ms.Called(ctx, courseID, lessonTitle, lessonText, lessonLanguage)

	var lesson *domain.Lesson
	if args.Get(0) == nil {
		lesson = nil
	} else {
		lesson = args.Get(0).(*domain.Lesson)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return lesson, err
}

func (ms *mockService) CreateLesson(ctx context.Context, title, text, language string) (*domain.Lesson, *domain.Error) {
	args := ms.Called(ctx, title, text, language)

	var lesson *domain.Lesson
	if args.Get(0) == nil {
		lesson = nil
	} else {
		lesson = args.Get(0).(*domain.Lesson)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return lesson, err
}

func (ms *mockService) PatchCourse(ctx context.Context, courseID string, patch []byte) *domain.Error {
	args := ms.Called(ctx, courseID, patch)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

// type mockAuthTool struct {
// 	mock.Mock
// }

// func (ma *mockAuthTool) GenerateTokenString(userID string) (string, error) {
// 	args := ma.Called(userID)

// 	return args.Get(0).(string), args.Error(1)
// }

// Creates a logger instance that discards all output
func nullLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	return logger
}

// errReader is intended to help us test - mainly in the corner case of error handling in the case of defective request/response bodies
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to read")
}

// type mockIDTool struct {
// 	mock.Mock
// }

// func (mit *mockIDTool) New() (string, error) {
// 	args := mit.Called()

// 	return args.Get(0).(string), args.Error(1)
// }
