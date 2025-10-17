package user

import (
	"time"

	param "github.com/NSObjects/go-template/internal/api/service/param"
	domain "github.com/NSObjects/go-template/internal/domain/user"
)

// AssembleListUsersQuery converts HTTP request DTO into domain query value object.
func AssembleListUsersQuery(req param.UserListUsersRequest) domain.ListUsersQuery {
	return domain.NewListUsersQuery(req.Page, req.Count)
}

// AssembleCreateUser converts create DTO into a domain aggregate.
func AssembleCreateUser(req param.UserCreateRequest) (domain.User, error) {
	username, err := domain.NewUsername(req.Username)
	if err != nil {
		return domain.User{}, err
	}
	email, err := domain.NewEmail(req.Email)
	if err != nil {
		return domain.User{}, err
	}
	age, err := domain.NewAge(req.Age)
	if err != nil {
		return domain.User{}, err
	}
	return domain.NewUserDraft(username, email, age), nil
}

// AssembleUpdateUser converts update DTO and identifier into a domain aggregate.
func AssembleUpdateUser(id int64, req param.UserUpdateRequest) (domain.User, error) {
	username, err := domain.NewUsername(req.Username)
	if err != nil {
		return domain.User{}, err
	}
	email, err := domain.NewEmail(req.Email)
	if err != nil {
		return domain.User{}, err
	}
	age, err := domain.NewAge(req.Age)
	if err != nil {
		return domain.User{}, err
	}
	draft := domain.NewUserDraft(username, email, age)
	return draft.WithID(domain.ID(id)), nil
}

// AssembleUserListResponse converts domain aggregates into response DTOs.
func AssembleUserListResponse(users []domain.User) []param.UserListItem {
	items := make([]param.UserListItem, 0, len(users))
	for _, u := range users {
		items = append(items, param.UserListItem{
			Id:        u.ID.Int64(),
			Username:  u.Username.String(),
			Email:     u.Email.String(),
			Age:       u.Age.Int(),
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	return items
}

// AssembleUserDataResponse converts a domain aggregate into a detail response DTO.
func AssembleUserDataResponse(u domain.User) param.UserData {
	return param.UserData{
		Id:        u.ID.Int64(),
		Username:  u.Username.String(),
		Email:     u.Email.String(),
		Age:       u.Age.Int(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// AssembleUserID converts primitive identifier into domain ID value object.
func AssembleUserID(id int64) domain.ID {
	return domain.ID(id)
}

// AssembleUserFromTimestamps reconstructs a domain aggregate with persisted timestamps.
func AssembleUserFromTimestamps(id int64, username, email string, age int, createdAt, updatedAt time.Time) (domain.User, error) {
	uname, err := domain.NewUsername(username)
	if err != nil {
		return domain.User{}, err
	}
	mail, err := domain.NewEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	years, err := domain.NewAge(age)
	if err != nil {
		return domain.User{}, err
	}
	draft := domain.NewUserDraft(uname, mail, years).WithID(domain.ID(id))
	return draft.WithTimestamps(createdAt, updatedAt), nil
}
