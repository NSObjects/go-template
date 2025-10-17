package user

// ListUsersQuery defines pagination query for listing users.
type ListUsersQuery struct {
	page  int
	count int
}

// NewListUsersQuery normalises pagination values.
func NewListUsersQuery(page, count int) ListUsersQuery {
	if page <= 0 {
		page = 1
	}
	if count <= 0 {
		count = 10
	}
	return ListUsersQuery{page: page, count: count}
}

func (q ListUsersQuery) Page() int {
	return q.page
}

func (q ListUsersQuery) Size() int {
	return q.count
}

func (q ListUsersQuery) Offset() int {
	return (q.page - 1) * q.count
}

func (q ListUsersQuery) Limit() int {
	return q.count
}
