package user

import (
	"time"

	"github.com/NSObjects/go-template/internal/api/data/model"
	domain "github.com/NSObjects/go-template/internal/domain/user"
)

// AssembleDomainUser builds a domain aggregate from persistence model.
func AssembleDomainUser(m model.User) (domain.User, error) {
	return AssembleUserFromTimestamps(m.ID, m.Username, m.Email, int(m.Age), m.CreatedAt, m.UpdatedAt)
}

// AssembleModelUserForCreate prepares a persistence model for create operations.
func AssembleModelUserForCreate(u domain.User, now time.Time) model.User {
	return model.User{
		Username:  u.Username.String(),
		Email:     u.Email.String(),
		Age:       int32(u.Age.Int()),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AssembleModelUserForUpdate prepares a persistence model for update operations.
func AssembleModelUserForUpdate(u domain.User, now time.Time) model.User {
	return model.User{
		ID:        u.ID.Int64(),
		Username:  u.Username.String(),
		Email:     u.Email.String(),
		Age:       int32(u.Age.Int()),
		UpdatedAt: now,
	}
}
