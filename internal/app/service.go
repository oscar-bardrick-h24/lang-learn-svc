package app

import (
	"context"

	"github.com/OJOMB/lang-learn-svc/internal/pkg/domain"
)

type Service interface {
	// Users
	CreateUser(ctx context.Context, email, pword, fname, lname, profPic string) (*domain.User, *domain.Error)
	GetUser(ctx context.Context, userID string) (*domain.User, *domain.Error)
	GetAuthenticatedUser(ctx context.Context, email, password string) (*domain.User, *domain.Error)
	GetUserCourses(ctx context.Context, userID string) ([]domain.UserCourse, *domain.Error)
	PatchUser(ctx context.Context, userID string, patch []byte) *domain.Error
	SetUserProfilePic(ctx context.Context, userID, profPic string) *domain.Error
	SetUserPassword(ctx context.Context, userID, password string) *domain.Error
	DeleteUser(ctx context.Context, userID string) *domain.Error
	EnrollUser(ctx context.Context, userID, courseID string) *domain.Error

	// Languages
	CreateLanguage(ctx context.Context, name, code string) *domain.Error
	GetLanguages(ctx context.Context) ([]domain.Language, *domain.Error)
	UpdateLanguage(ctx context.Context, code, name string) *domain.Error
	DeleteLanguage(ctx context.Context, code string) *domain.Error
	DeleteLanguages(ctx context.Context) *domain.Error

	// Courses
	CreateCourse(ctx context.Context, title string, lessons []domain.Lesson) (*domain.Course, *domain.Error)
	GetCourse(ctx context.Context, courseID string) (*domain.Course, *domain.Error)
	GetCourses(ctx context.Context, createdByID string) ([]domain.Course, *domain.Error)
	PatchCourse(ctx context.Context, courseID string, patch []byte) *domain.Error
	DeleteCourse(ctx context.Context, courseID string) *domain.Error
	AppendNewLessonToCourse(ctx context.Context, courseID, lessonTitle, lessonText, lessonLanguage string) (*domain.Lesson, *domain.Error)

	// Lessons
	CreateLesson(ctx context.Context, title, text, language string) (*domain.Lesson, *domain.Error)
	GetLessons(ctx context.Context) ([]domain.Lesson, *domain.Error)
	GetLesson(ctx context.Context, lessonID string) (*domain.Lesson, *domain.Error)
	DeleteLesson(ctx context.Context, lessonID string) *domain.Error
	UpdateLesson(ctx context.Context, lessonID, title, text, language string) *domain.Error
}
