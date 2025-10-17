package user

import (
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/api/data/model"
	param "github.com/NSObjects/go-template/internal/api/service/param"
	domain "github.com/NSObjects/go-template/internal/domain/user"
)

func TestAssembleCreateUser(t *testing.T) {
	req := param.UserCreateRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Age:      18,
	}

	aggregate, err := AssembleCreateUser(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if aggregate.Username.String() != req.Username {
		t.Fatalf("expected username %s got %s", req.Username, aggregate.Username.String())
	}
	if aggregate.Email.String() != req.Email {
		t.Fatalf("expected email %s got %s", req.Email, aggregate.Email.String())
	}
	if aggregate.Age.Int() != req.Age {
		t.Fatalf("expected age %d got %d", req.Age, aggregate.Age.Int())
	}
}

func TestAssembleCreateUser_Invalid(t *testing.T) {
	_, err := AssembleCreateUser(param.UserCreateRequest{})
	if err == nil {
		t.Fatal("expected error for invalid data")
	}
	if err != domain.ErrInvalidUsername {
		t.Fatalf("expected invalid username error got %v", err)
	}
}

func TestAssembleUserListResponse(t *testing.T) {
	username, _ := domain.NewUsername("bob")
	email, _ := domain.NewEmail("bob@example.com")
	age, _ := domain.NewAge(20)
	now := time.Now()
	user := domain.NewUser(domain.ID(1), username, email, age, now, now)

	resp := AssembleUserListResponse([]domain.User{user})
	if len(resp) != 1 {
		t.Fatalf("expected list length 1 got %d", len(resp))
	}
	if resp[0].Id != 1 {
		t.Fatalf("expected id 1 got %d", resp[0].Id)
	}
	if resp[0].Username != username.String() {
		t.Fatalf("expected username %s got %s", username.String(), resp[0].Username)
	}
	if resp[0].Email != email.String() {
		t.Fatalf("expected email %s got %s", email.String(), resp[0].Email)
	}
	if resp[0].Age != age.Int() {
		t.Fatalf("expected age %d got %d", age.Int(), resp[0].Age)
	}
}

func TestAssembleDomainUser(t *testing.T) {
	now := time.Now()
	modelUser := model.User{
		ID:        3,
		Username:  "cathy",
		Email:     "cathy@example.com",
		Age:       22,
		CreatedAt: now,
		UpdatedAt: now,
	}

	aggregate, err := AssembleDomainUser(modelUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if aggregate.ID.Int64() != modelUser.ID {
		t.Fatalf("expected id %d got %d", modelUser.ID, aggregate.ID.Int64())
	}
	if aggregate.Username.String() != modelUser.Username {
		t.Fatalf("expected username %s got %s", modelUser.Username, aggregate.Username.String())
	}
	if aggregate.Email.String() != modelUser.Email {
		t.Fatalf("expected email %s got %s", modelUser.Email, aggregate.Email.String())
	}
	if aggregate.Age.Int() != int(modelUser.Age) {
		t.Fatalf("expected age %d got %d", modelUser.Age, aggregate.Age.Int())
	}
}

func TestAssembleModelUserForCreate(t *testing.T) {
	username, _ := domain.NewUsername("demo")
	email, _ := domain.NewEmail("demo@example.com")
	age, _ := domain.NewAge(30)
	draft := domain.NewUserDraft(username, email, age)
	now := time.Now()

	record := AssembleModelUserForCreate(draft, now)
	if record.CreatedAt != now {
		t.Fatalf("expected created at %v got %v", now, record.CreatedAt)
	}
	if record.Username != username.String() {
		t.Fatalf("expected username %s got %s", username.String(), record.Username)
	}
	if int(record.Age) != age.Int() {
		t.Fatalf("expected age %d got %d", age.Int(), record.Age)
	}
}

func TestAssembleModelUserForUpdate(t *testing.T) {
	username, _ := domain.NewUsername("demo")
	email, _ := domain.NewEmail("demo@example.com")
	age, _ := domain.NewAge(30)
	now := time.Now()
	aggregate := domain.NewUser(domain.ID(5), username, email, age, now, now)

	record := AssembleModelUserForUpdate(aggregate, now)
	if record.ID != 5 {
		t.Fatalf("expected id 5 got %d", record.ID)
	}
	if record.Username != username.String() {
		t.Fatalf("expected username %s got %s", username.String(), record.Username)
	}
	if record.UpdatedAt != now {
		t.Fatalf("expected updated at %v got %v", now, record.UpdatedAt)
	}
}
