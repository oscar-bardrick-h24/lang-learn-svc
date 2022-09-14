package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/contextual"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

///////////////////
//  CreateUser  //
/////////////////

func TestCreateUser_successPath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	var (
		uID        = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		firstName  = "John"
		lastName   = "Doe"
		email      = "test@example.com"
		password   = "password"
		profilePic = "https://s3.eu-west-1.amazonaws.com/bucket/object"
		saltedHash = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"
	)

	expectedUser := User{
		ID:         uID,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		ProfilePic: profilePic,
		Password:   saltedHash,
	}

	mIDt.On("New").Return(uID, nil).Once()
	mIDt.On("IsValid", uID).Return(true).Once()

	mpt.On("New", password).Return(saltedHash, nil).Once()

	mr.On("CreateUser", mock.Anything, expectedUser).Return(nil).Once()

	svc := NewService(nullLogger(), mr, mIDt, mpt)
	user, err := svc.CreateUser(context.Background(), email, password, firstName, lastName, profilePic)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedUser, *user)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
	mpt.AssertExpectations(t)
}

///////////////
// PatchUser //
///////////////

func TestPatchUser_successPath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	var (
		ctx = contextual.SetSubjectID(context.Background(), "9abc46be-3bcd-42b1-aeb2-ac6ff557a580")

		uID        = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		firstName  = "John"
		lastName   = "Doe"
		email      = "test@example.com"
		profilePic = "https://s3.eu-west-1.amazonaws.com/bucket/object"
		saltedHash = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"

		newFirstName = "Jane"
		newLastName  = "Poe"
		newEmail     = "test1@example.com"
	)

	originalUser := &User{
		ID:         uID,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		ProfilePic: profilePic,
		Password:   saltedHash,
	}

	patchedUser := User{
		ID:         uID,
		Email:      newEmail,
		FirstName:  newFirstName,
		LastName:   newLastName,
		ProfilePic: profilePic,
		Password:   saltedHash,
	}

	patchJSON := []byte(fmt.Sprintf(`[
		{ "op": "replace", "path": "/first_name", "value": "%s" },
		{ "op": "replace", "path": "/last_name", "value": "%s" },
		{ "op": "replace", "path": "/email", "value": "%s" }
	]`, newFirstName, newLastName, newEmail))

	mIDt.On("IsValid", uID).Return(true).Twice()

	mr.On("GetUser", mock.Anything, originalUser.ID).Return(originalUser, nil).Once()
	mr.On("UpdateUser", mock.Anything, patchedUser).Return(nil).Once()

	svc := NewService(nullLogger(), mr, mIDt, mpt)
	err := svc.PatchUser(ctx, originalUser.ID, patchJSON)
	assert.Nil(t, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
	mpt.AssertExpectations(t)
}

func TestPatchUser_attemptToPatchUnpatchableAttributes_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	var (
		ctx = contextual.SetSubjectID(context.Background(), "9abc46be-3bcd-42b1-aeb2-ac6ff557a580")

		uID        = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		firstName  = "John"
		lastName   = "Doe"
		email      = "test@example.com"
		profilePic = "https://s3.eu-west-1.amazonaws.com/bucket/object"
		saltedHash = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"
	)

	originalUser := &User{
		ID:         uID,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		ProfilePic: profilePic,
		Password:   saltedHash,
	}

	patchJSON := []byte(`[
		{ "op": "replace", "path": "/id", "value": "new-id" },
		{ "op": "replace", "path": "/profile_pic", "value": "https://s3.eu-west-1.amazonaws.com/bucket/new-object" },
		{ "op": "replace", "path": "/password", "value": "newpassword" }
	]`)

	mIDt.On("IsValid", uID).Return(true).Once()

	mr.On("GetUser", mock.Anything, originalUser.ID).Return(originalUser, nil).Once()
	// mr.On("UpdateUser", mock.Anything, patchedUser).Return(nil).Once()

	svc := NewService(nullLogger(), mr, mIDt, mpt)
	err := svc.PatchUser(ctx, originalUser.ID, patchJSON)
	assert.Nil(t, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
	mpt.AssertExpectations(t)
}
