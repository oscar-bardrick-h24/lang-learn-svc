package domain

import (
	"fmt"
	"net/mail"
	"time"
)

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ProfilePic string    `json:"profile_pic"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type userPatchableAttributes struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewUser(id, email, password, fname, lname, profPic string) *User {
	return &User{
		ID:         id,
		Email:      email,
		FirstName:  fname,
		LastName:   lname,
		ProfilePic: profPic,
		Password:   password,
	}
}

func (u *User) getPatchableAttributes() userPatchableAttributes {
	return userPatchableAttributes{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (u *User) patchAttributes(p userPatchableAttributes) error {
	if p.Email == u.Email && p.FirstName == u.FirstName && p.LastName == u.LastName {
		return fmt.Errorf("patch effects no change")
	}

	u.Email = p.Email
	u.FirstName = p.FirstName
	u.LastName = p.LastName

	return nil
}

func (u *User) Validate(idValidator IDValidator, pwordValidator PasswordValidator) error {
	if !idValidator.IsValid(u.ID) {
		return fmt.Errorf("id format is invalid")
	}

	// Email
	switch {
	case u.Email == "":
		return fmt.Errorf("email must not be empty")
	case len(u.Email) > 255:
		return fmt.Errorf("email must not be longer than 255 characters")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("email format is invalid: %v", err)
	}

	// Password
	switch {
	case u.Password == "":
		return fmt.Errorf("password must not be empty")
	case len(u.Password) > 255:
		return fmt.Errorf("password must not be longer than 255 characters")
	}

	// if !pwordValidator.IsValid(u.Password) {
	// 	return fmt.Errorf("password is not in expected salted hash format")
	// }

	if u.FirstName == "" {
		return fmt.Errorf("first name must not be empty")
	}

	if u.LastName == "" {
		return fmt.Errorf("last name must not be empty")
	}

	// TODO: validate profile pic

	return nil
}
