package domain

import "context"

type Repo interface {
	// Users
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserCourses(ctx context.Context, userID string) ([]UserCourse, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, userID string) error
	EnrollUser(ctx context.Context, userID, courseID string) error

	// Languages
	CreateLanguage(ctx context.Context, language Language) error
	GetLanguages(ctx context.Context) ([]Language, error)
	UpdateLanguage(ctx context.Context, lang Language) error
	DeleteLanguage(ctx context.Context, code string) error
	DeleteLanguages(ctx context.Context) error

	// Courses
	CreateCourse(ctx context.Context, course Course, lessons ...Lesson) error
	GetCourse(ctx context.Context, courseID string) (*Course, error)
	GetCourses(ctx context.Context) ([]Course, error)
	GetCoursesByCreator(ctx context.Context, createdBy string) ([]Course, error)
	UpdateCourse(ctx context.Context, course Course) error
	DeleteCourse(ctx context.Context, courseID string) error
	AppendNewLessonToCourse(ctx context.Context, courseID string, lesson Lesson) error

	// Lessons
	CreateLesson(ctx context.Context, lesson Lesson) error
	GetLesson(ctx context.Context, lessonID string) (*Lesson, error)
	GetLessons(ctx context.Context) ([]Lesson, error)
	DeleteLesson(ctx context.Context, lessonID string) error
	UpdateLesson(ctx context.Context, lesson Lesson) error
}
