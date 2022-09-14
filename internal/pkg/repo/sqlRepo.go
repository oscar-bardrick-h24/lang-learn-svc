package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
	"github.com/sirupsen/logrus"
)

const componentSQLRepo = "SQLRepo"

type SQLRepo struct {
	db     *sql.DB
	logger *logrus.Entry
}

func NewSQLRepo(db *sql.DB, logger *logrus.Logger) *SQLRepo {
	return &SQLRepo{
		db:     db,
		logger: logger.WithField("component", componentSQLRepo),
	}
}

func (r *SQLRepo) Ping() error {
	return r.db.Ping()
}

///////////
// Users //
///////////

// CreateUser creates a new user record in the database
func (r *SQLRepo) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, email, password, first_name, last_name, profile_pic) VALUES ($1, $2, $3, $4, $5, $6);`,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.ProfilePic,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateUser").Error("failed to create user")
		return err
	}

	return nil
}

// GetUser retrieves a user record from the database by ID
func (r *SQLRepo) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password, first_name, last_name, profile_pic, created_at, updated_at
		FROM users
		WHERE id = $1;`,
		userID,
	).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.ProfilePic,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetUser").Error("failed to get user")
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user record from the database by email address
func (r *SQLRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password, first_name, last_name, profile_pic, created_at, updated_at
		FROM users WHERE email = $1;`,
		email,
	).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetUserByEmail").Error("failed to get user")
		return nil, err
	}

	return &user, nil
}

// GetUserCourses retrieves all user_courses records for a given user
func (r *SQLRepo) GetUserCourses(ctx context.Context, userID string) ([]domain.UserCourse, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT user_id, course_id, active_lesson_id, created_at, updated_at FROM user_courses WHERE user_id = $1;`,
		userID,
	)
	if err != nil && err == sql.ErrNoRows {
		return []domain.UserCourse{}, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetUserCourses").Error("failed to get user courses")
		return nil, err
	}

	defer rows.Close()

	userCourses := make([]domain.UserCourse, 0)
	for rows.Next() {
		var userCourse domain.UserCourse
		if err = rows.Scan(
			&userCourse.UserID, &userCourse.CourseID, &userCourse.ActiveLessonID, &userCourse.CreatedAt, &userCourse.UpdatedAt,
		); err != nil {
			r.logger.WithError(err).WithField("method", "GetUserCourses").Error("failed to scan row")
			return nil, err
		}

		userCourses = append(userCourses, userCourse)
	}

	return userCourses, nil
}

// UpdateUser updates a user record in the database
func (r *SQLRepo) UpdateUser(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET email = $1, password = $2, first_name = $3, last_name = $4, profile_pic = $5, updated_at = current_timestamp 
		WHERE id = $6;`,
		user.Email, user.Password, user.FirstName, user.LastName, user.ProfilePic, user.ID,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "UpdateUser").Error("failed to update user")
		return err
	}

	return nil
}

// DeleteUser deletes a user from the database
func (r *SQLRepo) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM users WHERE id = $1;`,
		userID,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteUser").Error("failed to delete user")
		return err
	}

	return nil
}

// EnrollUser creates a user_course record to represent the fact that a user is enrolled in a course
func (r *SQLRepo) EnrollUser(ctx context.Context, userID, courseID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user_courses (user_id, course_id) VALUES ($1, $2);`,
		userID, courseID,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "EnrollUser").Error("failed to enroll user")
		return err
	}

	return nil
}

///////////////
// Languages //
///////////////

// CreateLanguage creates a language in the database
func (r *SQLRepo) CreateLanguage(ctx context.Context, language domain.Language) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO languages (code, name) VALUES ($1, $2);`,
		language.Code, language.Name,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateLanguage").Error("failed to create language")
		return err
	}

	return nil
}

// GetLanguages retrieves all languages from the database
func (r *SQLRepo) GetLanguages(ctx context.Context) ([]domain.Language, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT code, name, created_at, updated_at FROM languages;`,
	)
	if err != nil && err == sql.ErrNoRows {
		return []domain.Language{}, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetLanguages").Error("failed to get languages")
		return nil, err
	}

	defer rows.Close()

	var languages []domain.Language
	for rows.Next() {
		var language domain.Language
		if err := rows.Scan(&language.Code, &language.Name, &language.CreatedAt, &language.UpdatedAt); err != nil {
			r.logger.WithError(err).WithField("method", "GetLanguages").Error("failed to scan language")
			return nil, err
		}

		languages = append(languages, language)
	}

	return languages, nil
}

// UpdateLanguage updates a language in the database with the given code
func (r *SQLRepo) UpdateLanguage(ctx context.Context, language domain.Language) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE languages SET name = $1, updated_at = current_timestamp WHERE code = $2;`,
		language.Name, language.Code,
	)
	if err != nil && err != sql.ErrNoRows {
		r.logger.WithError(err).WithField("method", "UpdateLanguage").Error("failed to update language")
	}

	return err
}

// DeleteLanguage deletes a language from the database by code
func (r *SQLRepo) DeleteLanguage(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM languages WHERE code = $1;`,
		code,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteLanguage").Error("failed to delete language")
		return err
	}

	return nil
}

// DeleteLanguages deletes all languages from the database
func (r *SQLRepo) DeleteLanguages(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM languages`,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteLanguages").Error("failed to delete languages")
		return err
	}

	return nil
}

/////////////
// Courses //
/////////////

// CreateCourse creates a course record in the database as well as optionally all given lessons in a transaction
func (r *SQLRepo) CreateCourse(ctx context.Context, course domain.Course, lessons ...domain.Lesson) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateCourse").Error("failed to begin transaction")
		return err
	}

	// handle lesson creation
	if err := r.createLessonsTx(ctx, tx, lessons...); err != nil {
		r.logger.WithError(err).WithField("method", "CreateCourse").Error("failed to create lessons for course")
		tx.Rollback()
		return err
	}

	lessonIDs, err := json.Marshal(course.LessonIDs)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateCourse").Error("failed to marshal lesson ids")
		tx.Rollback()
		return err
	}

	// handle course creation
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO courses (id, title, lesson_ids, created_by) VALUES ($1, $2, $3, $4);`,
		course.ID, course.Title, lessonIDs, course.CreatedBy,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateCourse").Error("failed to create course")
		tx.Rollback()
		return err
	}

	// associate lessons with course if they exist
	if err := r.createCourseLessonsTx(ctx, tx, course); err != nil {
		r.logger.WithError(err).WithField("method", "CreateCourse").Error("failed to create course lessons records")
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).WithField("method", "CreateCourse").Error("failed to commit transaction")
		tx.Rollback()
		return err
	}

	return nil
}

// GetCourse retrieves a course record from the database by courseID
func (r *SQLRepo) GetCourse(ctx context.Context, courseID string) (*domain.Course, error) {
	var course domain.Course
	lessonIDsBytes := make([]byte, 0)
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, lesson_ids, created_by, created_at, updated_at FROM courses WHERE id = $1;`,
		courseID,
	).Scan(&course.ID, &course.Title, &lessonIDsBytes, &course.CreatedBy, &course.CreatedAt, &course.UpdatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetCourse").Error("failed to get course")
		return nil, err
	}

	if err := json.Unmarshal(lessonIDsBytes, &course.LessonIDs); err != nil {
		r.logger.WithError(err).WithField("method", "GetCourse").Error("failed to unmarshal lesson ids")
		return nil, err
	}

	return &course, nil
}

// GetCourses retrieves all courses from the database
func (r *SQLRepo) GetCourses(ctx context.Context) ([]domain.Course, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, lesson_ids, created_by, created_at, updated_at FROM courses;`,
	)
	if err != nil && err == sql.ErrNoRows {
		return []domain.Course{}, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetCourses").Error("failed to get courses")
		return nil, err
	}

	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var course domain.Course
		if err := rows.Scan(
			&course.ID, &course.Title, &course.LessonIDs, &course.CreatedBy, &course.CreatedAt, &course.UpdatedAt,
		); err != nil {
			r.logger.WithError(err).WithField("method", "GetLanguages").Error("failed to scan language")
			return nil, err
		}

		courses = append(courses, course)
	}

	return courses, nil
}

// GetCoursesByCreator retrieves all courses from the database created by the given user
func (r *SQLRepo) GetCoursesByCreator(ctx context.Context, createdBy string) ([]domain.Course, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, lesson_ids, created_by, created_at, updated_at FROM courses WHERE created_by = $1;`,
		createdBy,
	)
	if err != nil && err == sql.ErrNoRows {
		return []domain.Course{}, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetCoursesByCreator").Error("failed to get courses by creator")
		return nil, err
	}

	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var course domain.Course
		if err := rows.Scan(
			&course.ID, &course.Title, &course.LessonIDs, &course.CreatedBy, &course.CreatedAt, &course.UpdatedAt,
		); err != nil {
			r.logger.WithError(err).WithField("method", "GetCoursesByCreator").Error("failed to scan course")
			return nil, err
		}
	}

	return courses, nil
}

func (r *SQLRepo) DeleteCourse(ctx context.Context, courseID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteCourse").Error("failed to begin transaction")
		return err
	}

	// delete course
	if _, err := tx.ExecContext(ctx, `DELETE FROM courses WHERE id = $1;`, courseID); err != nil && err == sql.ErrNoRows {
		tx.Rollback()
		return nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteCourse").Errorf("failed to delete course %s", courseID)
		tx.Rollback()
		return err
	}

	// delete course_lessons records
	if _, err := tx.ExecContext(ctx, `DELETE FROM course_lessons WHERE course_id = $1;`, courseID); err != nil && err != sql.ErrNoRows {
		r.logger.WithError(err).WithField("method", "DeleteCourse").
			Errorf("failed to delete course_lessons records for course %s", courseID)
		tx.Rollback()
		return err
	}

	// delete user_courses records
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_courses WHERE course_id = $1;`, courseID); err != nil && err != sql.ErrNoRows {
		r.logger.WithError(err).WithField("method", "DeleteCourse").
			Errorf("failed to delete user_courses records for course %s", courseID)
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).WithField("method", "DeleteCourse").Error("failed to commit transaction")
		tx.Rollback()
		return err
	}

	return nil
}

// AppendNewLessonToCourse adds the new lesson to the course at the end of the lesson list
func (r *SQLRepo) AppendNewLessonToCourse(ctx context.Context, courseID string, lesson domain.Lesson) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).WithField("method", "AddNewLessonToCourse").Error("failed to begin transaction")
		return err
	}

	// get course
	course, err := r.GetCourse(ctx, courseID)
	if err != nil {
		r.logger.WithError(err).WithField("method", "AddNewLessonToCourse").Error("failed to get course")
		tx.Rollback()
		return err
	}

	// add lesson to course
	course.LessonIDs = append(course.LessonIDs, lesson.ID)

	// update course
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE courses SET lesson_ids = $1, updated_at = current_timestamp WHERE id = $2;`,
		course.LessonIDs, course.ID,
	); err != nil {
		r.logger.WithError(err).WithField("method", "AddNewLessonToCourse").Error("failed to update course")
		tx.Rollback()
		return err
	}

	// insert lesson
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO lessons (id, title, text, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, current_timestamp);`,
		lesson.ID, lesson.Title, lesson.Text, lesson.CreatedBy, lesson.CreatedAt, lesson.UpdatedAt,
	); err != nil {
		r.logger.WithError(err).WithField("method", "AddNewLessonToCourse").Error("failed to insert lesson")
		tx.Rollback()
		return err
	}

	// insert course_lessons record
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO course_lessons (course_id, lesson_id) VALUES ($1, $2);`,
		course.ID, lesson.ID,
	); err != nil {
		r.logger.WithError(err).WithField("method", "AddNewLessonToCourse").Error("failed to insert course_lessons record")
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).WithField("method", "AddNewLessonToCourse").Error("failed to commit transaction")
		tx.Rollback()
		return err
	}

	return nil
}

func (r *SQLRepo) UpdateCourse(ctx context.Context, course domain.Course) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).WithField("method", "UpdateCourse").Error("failed to begin transaction")
		return err
	}

	// update course
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE courses SET title = $1, lesson_ids = $2, updated_at = current_timestamp WHERE id = $3;`,
		course.Title, course.LessonIDs, course.ID,
	); err != nil {
		r.logger.WithError(err).WithField("method", "UpdateCourse").Error("failed to update course")
		tx.Rollback()
		return err
	}

	// delete course_lessons records
	if _, err := tx.ExecContext(ctx, `DELETE FROM course_lessons WHERE course_id = $1;`, course.ID); err != nil && err != sql.ErrNoRows {
		r.logger.WithError(err).WithField("method", "UpdateCourse").
			Errorf("failed to delete course_lessons records for course %s", course.ID)
		tx.Rollback()
		return err
	}

	// insert course_lessons records
	if err := r.createCourseLessonsTx(ctx, tx, course); err != nil {
		r.logger.WithError(err).WithField("method", "UpdateCourse").
			Errorf("failed to insert course_lessons records for course %s", course.ID)
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).WithField("method", "UpdateCourse").Error("failed to commit transaction")
		tx.Rollback()
		return err
	}

	return nil
}

///////////////////
// CourseLessons //
///////////////////

// createCourseLessonsTx creates new course_lessons records in the database in a single query.
// course_lessons records associate lessons with a course and specifies the order of the lessons.
func (r *SQLRepo) createCourseLessonsTx(ctx context.Context, tx *sql.Tx, course domain.Course) error {
	if course.LessonIDs == nil || len(course.LessonIDs) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(course.LessonIDs)*2)
	query := `INSERT INTO course_lessons (course_id, lesson_id, lesson_number) VALUES `

	counter := 1
	for i, lessonID := range course.LessonIDs {
		query += fmt.Sprintf("($%d, $%d, $%d)", counter, counter+1, counter+2)
		if i != len(course.LessonIDs)-1 {
			query += `, `
		} else {
			query += `;`
		}

		args = append(args, course.ID, lessonID, i+1)
		counter += 3
	}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// getCourseLessonsTx returns the course_lessons records by course ID
// caller is responsible for committing or rolling back the transaction
func (r *SQLRepo) getCourseLessonsTx(ctx context.Context, tx *sql.Tx, courseID string) ([]string, error) {
	rows, err := tx.QueryContext(
		ctx,
		`SELECT lesson_id FROM course_lessons WHERE course_id = $1 ORDER BY lesson_number;`,
		courseID,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "GetCourseLessons").Error("failed to query course_lessons records")
		return nil, err
	}
	defer rows.Close()

	lessonIDs := make([]string, 0)
	for rows.Next() {
		var lessonID string
		if err := rows.Scan(&lessonID); err != nil {
			r.logger.WithError(err).WithField("method", "GetCourseLessons").Error("failed to scan course_lessons record")
			return nil, err
		}

		lessonIDs = append(lessonIDs, lessonID)
	}

	if err := rows.Err(); err != nil {
		r.logger.WithError(err).WithField("method", "GetCourseLessons").Error("failed to iterate over course_lessons records")
		return nil, err
	}

	return lessonIDs, nil
}

/////////////
// Lessons //
/////////////

// CreateLesson creates a lesson record in the database
func (r *SQLRepo) CreateLesson(ctx context.Context, lesson domain.Lesson) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO lessons (id, title, text, created_by) VALUES ($1, $2, $3, $4);`,
		lesson.ID, lesson.Title, lesson.Text, lesson.CreatedBy,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateLesson").Error("failed to create lesson")
		return err
	}

	return nil
}

// createLessons creates lessons in the database in a single query
// caller is responsible for commiting or rolling back the tx
func (r *SQLRepo) createLessonsTx(ctx context.Context, tx *sql.Tx, lessons ...domain.Lesson) error {
	if lessons == nil {
		return nil
	}

	args := make([]interface{}, 0, len(lessons)*5)
	query := `INSERT INTO lessons (id, title, text, language, created_by) VALUES `

	counter := 1
	for i, lesson := range lessons {
		query += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d)`, counter, counter+1, counter+2, counter+3, counter+4)
		if i != len(lessons)-1 {
			query += `, `
		} else {
			query += `;`
		}

		args = append(args, lesson.ID, lesson.Title, lesson.Text, lesson.Language, lesson.CreatedBy)
		counter += 5
	}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetLesson retreives a lesson from the database by id
func (r *SQLRepo) GetLesson(ctx context.Context, id string) (*domain.Lesson, error) {
	var lesson domain.Lesson
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, text, language, created_by, created_at, updated_at FROM lessons WHERE id = $1;`,
		id,
	).Scan(
		&lesson.ID, &lesson.Title, &lesson.Text, &lesson.Language, &lesson.CreatedBy, &lesson.CreatedAt, &lesson.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetLesson").Error("failed to get lesson")
		return nil, err
	}

	return &lesson, nil
}

// GetLessons retreives all lessons from the database
// returns empty slice if no lessons are found
func (r *SQLRepo) GetLessons(ctx context.Context) ([]domain.Lesson, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, text, language, created_by, created_at, updated_at FROM lessons;`,
	)
	if err != nil && err != sql.ErrNoRows {
		return []domain.Lesson{}, nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "GetLessons").Error("failed to get lessons")
		return nil, err
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		var lesson domain.Lesson
		if err := rows.Scan(
			&lesson.ID, &lesson.Title, &lesson.Text, &lesson.Language, &lesson.CreatedBy, &lesson.CreatedAt, &lesson.UpdatedAt,
		); err != nil {
			r.logger.WithError(err).WithField("method", "GetLessons").Error("failed to scan lesson")
			return nil, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// DeleteLesson deletes a lesson record from the database by id
func (r *SQLRepo) DeleteLesson(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteLesson").Error("failed to begin transaction")
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM lessons WHERE id = $1;`,
		id,
	); err != nil && err != sql.ErrNoRows {
		return nil
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "DeleteLesson").Error("failed to delete lesson")
		return err
	}

	// delete course_lessons records
	if _, err := tx.ExecContext(ctx, `DELETE FROM course_lessons WHERE lesson_id = $1;`, id); err != nil && err != sql.ErrNoRows {
		r.logger.WithError(err).WithField("method", "DeleteLesson").Errorf("failed to delete course_lessons records for lesson %s", id)
		tx.Rollback()
		return err
	}

	return nil
}

// UpdateLesson updates a lesson record in the database
func (r *SQLRepo) UpdateLesson(ctx context.Context, lesson domain.Lesson) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE lessons SET title = $1, text = $2, language = $3, updated_at = $4 WHERE id = $5;`,
		lesson.Title, lesson.Text, lesson.Language, lesson.UpdatedAt, lesson.ID,
	)
	if err != nil && err == sql.ErrNoRows {
		return err
	} else if err != nil {
		r.logger.WithError(err).WithField("method", "UpdateLesson").Error("failed to update lesson")
		return err
	}

	return nil
}
