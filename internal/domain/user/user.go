package user

import "time"

// User represents the aggregate root of the user domain.
type User struct {
	ID        ID
	Username  Username
	Email     Email
	Age       Age
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser constructs a User aggregate ensuring domain invariants are respected.
func NewUser(id ID, username Username, email Email, age Age, createdAt, updatedAt time.Time) User {
	return User{
		ID:        id,
		Username:  username,
		Email:     email,
		Age:       age,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// NewUserDraft creates a draft user aggregate for create operations.
func NewUserDraft(username Username, email Email, age Age) User {
	return User{
		Username: username,
		Email:    email,
		Age:      age,
	}
}

// WithID returns a new User with the provided identifier.
func (u User) WithID(id ID) User {
	u.ID = id
	return u
}

// WithTimestamps returns a new User with the provided timestamps.
func (u User) WithTimestamps(createdAt, updatedAt time.Time) User {
	u.CreatedAt = createdAt
	u.UpdatedAt = updatedAt
	return u
}
