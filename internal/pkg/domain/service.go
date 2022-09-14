package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/contextual"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sirupsen/logrus"
)

const (
	componentService = "Service"
)

type Service struct {
	logger *logrus.Entry
	repo   Repo

	passWordTool PasswordTool
	idTool       IDTool
}

func NewService(logger *logrus.Logger, repo Repo, idTool IDTool, passwordTool PasswordTool) *Service {
	return &Service{
		logger:       logger.WithField("component", componentService),
		repo:         repo,
		idTool:       idTool,
		passWordTool: passwordTool,
	}
}

///////////
// Users //
///////////

func (svc *Service) CreateUser(ctx context.Context, email, password, fname, lname, profPic string) (*User, *Error) {
	if email == "" || password == "" || fname == "" {
		return nil, newInvalidInputError("each of email, password, first_name must not be empty", nil)
	}

	id, err := svc.idTool.New()
	if err != nil {
		return nil, newSystemError("failed to generate valid ID", err)
	}

	// generate salted hash from plaintext password
	saltedHash, err := svc.passWordTool.New(password)
	if err != nil {
		return nil, newSystemError("failed to salt and hash password", err)
	}

	// create and validate user in memory
	user := NewUser(id, email, saltedHash, fname, lname, profPic)
	if err := user.Validate(svc.idTool, svc.passWordTool); err != nil {
		return nil, newInvalidInputError("user is invalid", err)
	}

	if err := svc.repo.CreateUser(ctx, *user); err != nil {
		return nil, newSystemError("failed to store new user", err)
	}

	return user, nil
}

func (svc *Service) GetUser(ctx context.Context, userID string) (*User, *Error) {
	if !svc.idTool.IsValid(userID) {
		return nil, newInvalidInputError("format of userID is invalid", nil)
	}

	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	}

	return user, nil
}

func (svc *Service) GetUserByEmail(ctx context.Context, email string) (*User, *Error) {
	if email == "" {
		return nil, newInvalidInputError("email must not be empty", nil)
	}

	user, err := svc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	}

	return user, nil
}

func (svc *Service) GetAuthenticatedUser(ctx context.Context, email, password string) (*User, *Error) {
	if email == "" {
		return nil, newInvalidInputError("email must not be empty", nil)
	} else if password == "" {
		return nil, newInvalidInputError("password must not be empty", nil)
	}

	user, err := svc.repo.GetUserByEmail(ctx, email)
	if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	} else if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	}

	valid, derr := svc.validateUserCredentials(ctx, user.Password, password)
	if derr != nil {
		return nil, derr
	}

	if !valid {
		return nil, newAuthorizationError("invalid credentials", nil)
	}

	return user, nil
}

// GetUserCourses returns all courses the user with the given ID is enrolled in
func (svc *Service) GetUserCourses(ctx context.Context, userID string) ([]UserCourse, *Error) {
	if !svc.idTool.IsValid(userID) {
		return nil, newInvalidInputError("format of userID is invalid", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != userID {
		return nil, newAuthorizationError("users cannot see which courses other users are enrolled in", nil)
	}

	userCourses, err := svc.repo.GetUserCourses(ctx, userID)
	if err != nil {
		return nil, newSystemError("failed to retrieve courses", err)
	}

	return userCourses, nil
}

// PatchUser updates the user with the given patch.
// Note that certain user fields cannot be patched: ID, Password, ProfilePic, CreatedAt, UpdatedAt.
// Attempts to patch these fields will be ignored
func (svc *Service) PatchUser(ctx context.Context, userID string, patchJSON []byte) *Error {
	if !svc.idTool.IsValid(userID) {
		return newInvalidInputError("format of userID is invalid", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != userID {
		return newAuthorizationError("users are only authorised to patch their own accounts", nil)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return newInvalidInputError("patch could not be decoded", err)
	}

	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return newResourceNotFoundError("user does not exist", nil)
	}

	currentUserJSON, err := json.Marshal(user.getPatchableAttributes())
	if err != nil {
		return newSystemError("failed to marshal existing user", err)
	}

	userPatch, dErr := svc.createUserPatch(currentUserJSON, patch)
	if dErr != nil {
		return dErr.WrapMessage("failed to patch user")
	} else if userPatch == nil {
		return nil
	}

	if user.patchAttributes(*userPatch) != nil {
		// patch does nothing so nothing to do - not an error according to RFC 5789
		return nil
	}

	// need to validate user post-patch to ensure we're not left in an invalid state
	if err := user.Validate(svc.idTool, svc.passWordTool); err != nil {
		return newInvalidInputError("patch would leave user in invalid state", err)
	}

	if err := svc.repo.UpdateUser(ctx, *user); err != nil {
		return newSystemError("failed to update user with patched attributes", err)
	}

	return nil
}

// createUserPatch creates a new user from the current user and the patch.
// https://jsonpatch.com/
func (svc *Service) createUserPatch(userAttr []byte, patch jsonpatch.Patch) (*userPatchableAttributes, *Error) {
	patchedUserAttr, err := patch.Apply(userAttr)
	if err != nil {
		return nil, newSystemError("failed to apply user patch", err)
	}

	// check if the patch actually changed anything
	if jsonpatch.Equal(userAttr, patchedUserAttr) {
		return nil, nil
	}

	var patchedUser userPatchableAttributes
	if err := json.Unmarshal(patchedUserAttr, &patchedUser); err != nil {
		return nil, newSystemError("could not unmarshal patched user", err)
	}

	return &patchedUser, nil
}

// DeleteUser deletes the user with the given ID
func (svc *Service) DeleteUser(ctx context.Context, userID string) *Error {
	if !svc.idTool.IsValid(userID) {
		return newInvalidInputError("format of userID is invalid", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != userID {
		return newAuthorizationError("users are only authorised to delete their own accounts", nil)
	}

	if err := svc.repo.DeleteUser(ctx, userID); err != nil {
		return newSystemError("failed to delete user", err)
	}

	return nil
}

// validateUserCredentials checks if the given credentials are valid.
// If they are, we return the User, if not we return an error
func (svc *Service) validateUserCredentials(ctx context.Context, storedPword, GivenPword string) (bool, *Error) {
	if err := svc.passWordTool.Check(storedPword, GivenPword); err != nil && strings.Contains(err.Error(), "mismatched hash and password") {
		return false, nil
	} else if err != nil {
		return false, newSystemError("failed to validate password", err)
	}

	return true, nil
}

// SetUserProfilePic sets the profile picture for the user with the given ID
func (svc *Service) SetUserProfilePic(ctx context.Context, userID string, profilePic string) *Error {
	if !svc.idTool.IsValid(userID) {
		return newInvalidInputError("format of userID is invalid", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != userID {
		return newAuthorizationError("users are only authorised to set their own profile picture", nil)
	}

	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return newResourceNotFoundError("user does not exist", nil)
	}

	user.ProfilePic = profilePic

	if err := svc.repo.UpdateUser(ctx, *user); err != nil {
		return newSystemError("failed to update user with profile pic", err)
	}

	return nil
}

// SetUserPassword sets the password for the user with the given ID
func (svc *Service) SetUserPassword(ctx context.Context, userID string, password string) *Error {
	if !svc.idTool.IsValid(userID) {
		return newInvalidInputError("format of userID is invalid", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != userID {
		return newAuthorizationError("users are only authorised to set their own password", nil)
	}

	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return newResourceNotFoundError("user does not exist", nil)
	}

	user.Password = password

	if err := svc.repo.UpdateUser(ctx, *user); err != nil {
		return newSystemError("failed to update user with password", err)
	}

	return nil
}

func (svc *Service) EnrollUser(ctx context.Context, userID, courseID string) *Error {
	if !svc.idTool.IsValid(userID) {
		return newInvalidInputError("format of userID is invalid", nil)
	}

	if !svc.idTool.IsValid(courseID) {
		return newInvalidInputError("format of courseID is invalid", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != userID {
		return newAuthorizationError("users are only authorised to enroll themselves on courses", nil)
	}

	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return newResourceNotFoundError("user does not exist", nil)
	}

	course, err := svc.repo.GetCourse(ctx, courseID)
	if err != nil {
		return newSystemError("failed to retrieve course", err)
	} else if course == nil {
		return newResourceNotFoundError("course does not exist", nil)
	}

	if err := svc.repo.EnrollUser(ctx, userID, courseID); err != nil {
		return newSystemError("failed to enroll user", err)
	}

	return nil
}

///////////////
// Languages //
///////////////

func (svc *Service) CreateLanguage(ctx context.Context, code, name string) *Error {
	if name == "" || code == "" {
		return newInvalidInputError("language name and code must not be empty", nil)
	}

	language := &Language{Code: code, Name: name}
	if err := svc.repo.CreateLanguage(ctx, *language); err != nil {
		return newSystemError("failed to create language", err)
	}

	return nil
}

func (svc *Service) GetLanguages(ctx context.Context) ([]Language, *Error) {
	languages, err := svc.repo.GetLanguages(ctx)
	if err != nil {
		return nil, newSystemError("failed to retrieve languages", err)
	}

	return languages, nil
}

// UpdateLanguage updates the language with the given code
func (svc *Service) UpdateLanguage(ctx context.Context, code, name string) *Error {
	if name == "" || code == "" {
		return newInvalidInputError("language name and code must not be empty", nil)
	}

	if err := svc.repo.UpdateLanguage(ctx, Language{Code: code, Name: name}); err != nil && err == sql.ErrNoRows {
		return newResourceNotFoundError("language does not exist", err)
	} else if err != nil {
		return newSystemError("failed to update language", err)
	}

	return nil
}

// DeleteLanguage deletes the language with the given code
func (svc *Service) DeleteLanguage(ctx context.Context, code string) *Error {
	// TODO: how should this be dealt with in terms of authorisation?
	// should probably add the concept of an admin user type to manage languages on the platform

	if code == "" {
		return newInvalidInputError("language code must not be empty", nil)
	}

	if err := svc.repo.DeleteLanguage(ctx, code); err != nil {
		return newSystemError("failed to delete language", err)
	}

	return nil
}

// DeleteLanguages deletes all languages
func (svc *Service) DeleteLanguages(ctx context.Context) *Error {
	// TODO: how should this be dealt with in terms of authorisation?
	// should probably add the concept of an admin user type to manage languages on the platform

	if err := svc.repo.DeleteLanguages(ctx); err != nil {
		return newSystemError("failed to delete languages", err)
	}

	return nil
}

/////////////
// Courses //
/////////////

// CreateCourse creates a new course and the accompanying lessons if necessary
func (svc *Service) CreateCourse(ctx context.Context, title string, lessons []Lesson) (*Course, *Error) {
	if title == "" {
		return nil, newInvalidInputError("course title must not be empty", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID == "" || !svc.idTool.IsValid(subjectID) {
		return nil, newSystemError("valid user authentication data not found", nil)
	}

	courseID, err := svc.idTool.New()
	if err != nil {
		return nil, newSystemError("failed to generate courseID", err)
	}

	course := Course{ID: courseID, Title: title, CreatedBy: subjectID, LessonIDs: make([]string, len(lessons))}

	for i := range lessons {
		lessonID, err := svc.idTool.New()
		if err != nil {
			return nil, newSystemError("failed to generate lessonID", err)
		}

		lessons[i].ID = lessonID
		lessons[i].CreatedBy = subjectID
		course.LessonIDs[i] = lessonID
	}

	if err := svc.repo.CreateCourse(ctx, course, lessons...); err != nil {
		return nil, newSystemError("failed to create course", err)
	}

	return &course, nil
}

// GetCourse retrieves the course with the given ID
func (svc *Service) GetCourse(ctx context.Context, courseID string) (*Course, *Error) {
	if !svc.idTool.IsValid(courseID) {
		return nil, newInvalidInputError("format of courseID is invalid", nil)
	}

	course, err := svc.repo.GetCourse(ctx, courseID)
	if err != nil {
		return nil, newSystemError("failed to retrieve course", err)
	} else if course == nil {
		return nil, newResourceNotFoundError("course does not exist", nil)
	}

	return course, nil
}

// GetCourses returns all courses, optionally filtered to those created by the given userID in createdBy
func (svc *Service) GetCourses(ctx context.Context, createdBy string) ([]Course, *Error) {
	var courses []Course
	var err error
	if createdBy == "" {
		courses, err = svc.repo.GetCourses(ctx)
		if err != nil {
			return nil, newSystemError("failed to retrieve courses", err)
		}
	} else {
		if !svc.idTool.IsValid(createdBy) {
			return nil, newInvalidInputError("format of userID given as createdBy is invalid", err)
		}

		courses, err = svc.repo.GetCoursesByCreator(ctx, createdBy)
		if err != nil {
			return nil, newSystemError("failed to retrieve courses", err)
		}
	}

	return courses, nil
}

func (svc *Service) DeleteCourse(ctx context.Context, courseID string) *Error {
	if !svc.idTool.IsValid(courseID) {
		return newInvalidInputError("format of courseID is invalid", nil)
	}

	course, err := svc.repo.GetCourse(ctx, courseID)
	if err != nil {
		return newSystemError("failed to retrieve course", err)
	} else if course == nil {
		return newResourceNotFoundError("course does not exist", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != course.CreatedBy {
		return newAuthorizationError("user is not authorised to delete this course", nil)
	}

	if err := svc.repo.DeleteCourse(ctx, courseID); err != nil {
		return newSystemError("failed to delete course", err)
	}

	return nil
}

// PatchCourse updates the course with the given patch.
// Note that certain course fields cannot be patched: ID, Password, ProfilePic, CreatedAt, UpdatedAt.
// Attempts to patch these fields will be ignored
func (svc *Service) PatchCourse(ctx context.Context, courseID string, patchJSON []byte) *Error {
	if !svc.idTool.IsValid(courseID) {
		return newInvalidInputError("format of courseID is invalid", nil)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return newInvalidInputError("patch could not be decoded", err)
	}

	course, err := svc.repo.GetCourse(ctx, courseID)
	if err != nil {
		return newSystemError("failed to retrieve course", err)
	} else if course == nil {
		return newResourceNotFoundError("course does not exist", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != course.CreatedBy {
		return newAuthorizationError("users are only authorised to patch their own courses", nil)
	}

	currentCourseJSON, err := json.Marshal(course.getPatchableAttributes())
	if err != nil {
		return newSystemError("failed to marshal existing user", err)
	}

	coursePatch, dErr := svc.createCoursePatch(currentCourseJSON, patch)
	if dErr != nil {
		return dErr.WrapMessage("failed to patch user")
	} else if coursePatch == nil {
		return nil
	}

	if course.patchAttributes(*coursePatch) != nil {
		// patch does nothing so nothing to do - not an error according to RFC 5789
		return nil
	}

	// need to validate course post-patch to ensure we're not left in an invalid state
	if err := course.Validate(svc.idTool, svc.passWordTool); err != nil {
		return newInvalidInputError("patch would leave user in invalid state", err)
	}

	if err := svc.repo.UpdateCourse(ctx, *course); err != nil {
		return newSystemError("failed to update user with patched attributes", err)
	}

	return nil
}

// createUserPatch creates a new user from the current user and the patch.
// https://jsonpatch.com/
func (svc *Service) createCoursePatch(userAttr []byte, patch jsonpatch.Patch) (*coursePatchableAttributes, *Error) {
	patchedCourseAttr, err := patch.Apply(userAttr)
	if err != nil {
		return nil, newSystemError("failed to apply course patch", err)
	}

	// check if the patch actually changed anything
	if jsonpatch.Equal(userAttr, patchedCourseAttr) {
		return nil, nil
	}

	var patchedCourse coursePatchableAttributes
	if err := json.Unmarshal(patchedCourseAttr, &patchedCourse); err != nil {
		return nil, newSystemError("could not unmarshal patched user", err)
	}

	return &patchedCourse, nil
}

func (svc *Service) AppendNewLessonToCourse(ctx context.Context, courseID string, lessonTitle, lessonText, lessonLang string) (*Lesson, *Error) {
	if lessonTitle == "" || lessonText == "" || lessonLang == "" {
		return nil, newInvalidInputError("lesson title, text and language must not be empty", nil)
	}

	if !svc.idTool.IsValid(courseID) {
		return nil, newInvalidInputError("format of courseID is invalid", nil)
	}

	course, err := svc.repo.GetCourse(ctx, courseID)
	if err != nil {
		return nil, newSystemError("failed to retrieve course", err)
	} else if course == nil {
		return nil, newResourceNotFoundError("course does not exist", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != course.CreatedBy {
		return nil, newAuthorizationError("user is not authorised to add a lesson to this course", nil)
	}

	lessonID, err := svc.idTool.New()
	if err != nil {
		return nil, newSystemError("failed to generate lessonID", err)
	}

	lesson := Lesson{
		ID:        lessonID,
		Title:     lessonTitle,
		Text:      lessonText,
		Language:  lessonLang,
		CreatedBy: subjectID,
	}
	if err := svc.repo.AppendNewLessonToCourse(ctx, courseID, lesson); err != nil {
		return nil, newSystemError("failed to add lesson to course", err)
	}

	return &lesson, nil
}

/////////////
// Lessons //
/////////////

// CreateLesson creates a new lesson record in the database.
func (svc *Service) CreateLesson(ctx context.Context, title, text, language string) (*Lesson, *Error) {
	if title == "" || text == "" || language == "" {
		return nil, newInvalidInputError("lesson title, text and language must not be empty", nil)
	}

	lessonID, err := svc.idTool.New()
	if err != nil {
		return nil, newSystemError("failed to generate lessonID", err)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID == "" || !svc.idTool.IsValid(subjectID) {
		return nil, newAuthorizationError("failed to retrieve valid authorised identity", nil)
	}

	lesson := Lesson{
		ID:        lessonID,
		Title:     title,
		Text:      text,
		Language:  language,
		CreatedBy: subjectID,
	}

	if err := svc.repo.CreateLesson(ctx, lesson); err != nil {
		return nil, newSystemError("failed to create lesson", err)
	}

	return &lesson, nil
}

// GetLesson retrieves the lesson record with the given ID.
func (svc *Service) GetLesson(ctx context.Context, lessonID string) (*Lesson, *Error) {
	if !svc.idTool.IsValid(lessonID) {
		return nil, newInvalidInputError("format of lessonID is invalid", nil)
	}

	lesson, err := svc.repo.GetLesson(ctx, lessonID)
	if err != nil {
		return nil, newSystemError("failed to retrieve lesson", err)
	} else if lesson == nil {
		return nil, newResourceNotFoundError("lesson does not exist", nil)
	}

	return lesson, nil
}

// GetLessons retrieves all lesson records
func (svc *Service) GetLessons(ctx context.Context) ([]Lesson, *Error) {
	lessons, err := svc.repo.GetLessons(ctx)
	if err != nil {
		return nil, newSystemError("failed to retrieve lessons", err)
	}

	return lessons, nil
}

func (svc *Service) DeleteLesson(ctx context.Context, lessonID string) *Error {
	if !svc.idTool.IsValid(lessonID) {
		return newInvalidInputError("format of lessonID is invalid", nil)
	}

	lesson, err := svc.repo.GetLesson(ctx, lessonID)
	if err != nil {
		return newSystemError("failed to delete lesson", err)
	} else if lesson == nil {
		return newResourceNotFoundError("lesson does not exist", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != lesson.CreatedBy {
		return newAuthorizationError("user is not authorized to delete this lesson", nil)
	}

	if err := svc.repo.DeleteLesson(ctx, lessonID); err != nil {
		return newSystemError("failed to delete lesson", err)
	}

	return nil
}

// UpdateLesson updates the lesson with the given data
func (svc *Service) UpdateLesson(ctx context.Context, lessonID, lessonTitle, lessonText, lessonLang string) *Error {
	if !svc.idTool.IsValid(lessonID) {
		return newInvalidInputError("format of lessonID is invalid", nil)
	}

	lesson, err := svc.repo.GetLesson(ctx, lessonID)
	if err != nil {
		return newSystemError("failed to retrieve lesson", err)
	} else if lesson == nil {
		return newResourceNotFoundError("lesson does not exist", nil)
	}

	subjectID := contextual.GetSubjectID(ctx)
	if subjectID != lesson.CreatedBy {
		return newAuthorizationError("user is not authorized to update this lesson", nil)
	}

	lesson.Title = lessonTitle
	lesson.Text = lessonText
	lesson.Language = lessonLang

	if err := svc.repo.UpdateLesson(ctx, *lesson); err != nil {
		return newSystemError("failed to update lesson", err)
	}

	return nil
}
