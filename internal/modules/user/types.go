package user

import "time"

// Query contains pagination input shared by user list requests.
type Query struct {
	Page  int `json:"page" form:"page" query:"page"`
	Count int `json:"count" form:"count" query:"count"`
}

// Limit returns the requested page size or the default page size.
func (q *Query) Limit() int {
	if q.Count == 0 {
		q.Count = 10
	}
	return q.Count
}

// Offset returns the zero-based row offset for the requested page.
func (q *Query) Offset() int {
	if q.Page > 0 {
		q.Page--
	}
	return q.Page * q.Count
}

// ListUsersRequest contains the user list query parameters.
type ListUsersRequest struct {
	Query
}

// CreateRequest contains the user creation payload.
type CreateRequest struct {
	Username string `param:"username" query:"username" form:"username" json:"username" xml:"username" validate:"required,min=3,max=20"`
	Email    string `param:"email" query:"email" form:"email" json:"email" xml:"email" validate:"required,email"`
	Age      int    `param:"age" query:"age" form:"age" json:"age" xml:"age" validate:"min=0,max=150"`
}

// UpdateRequest contains the user update payload.
type UpdateRequest struct {
	Username string `param:"username" query:"username" form:"username" json:"username" xml:"username" validate:"omitempty,min=3,max=20"`
	Email    string `param:"email" query:"email" form:"email" json:"email" xml:"email" validate:"omitempty,email"`
	Age      int    `param:"age" query:"age" form:"age" json:"age" xml:"age" validate:"omitempty,min=0,max=150"`
}

// ListItem is one user row returned by list operations.
type ListItem struct {
	Id        int64     `json:"id" validate:"required"`
	Username  string    `json:"username" validate:"required,min=3,max=20"`
	Email     string    `json:"email" validate:"required,email"`
	Age       int       `json:"age" validate:"min=0,max=150"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Data is the user detail response.
type Data struct {
	Username  string    `json:"username" validate:"required,min=3,max=20"`
	Email     string    `json:"email" validate:"required,email"`
	Age       int       `json:"age" validate:"min=0,max=150"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Id        int64     `json:"id" validate:"required"`
}
