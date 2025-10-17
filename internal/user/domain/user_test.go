package domain

import (
	"testing"
	"time"
)

func TestUser_Create(t *testing.T) {
	// 测试创建用户
	userID, err := NewUserID("user-123")
	if err != nil {
		t.Fatalf("Failed to create user ID: %v", err)
	}

	_, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	user, err := NewUser(userID, "testuser", "test@example.com", birthDate)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 验证用户属性
	if user.ID().Value() != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", user.ID().Value())
	}

	if user.Username() != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username())
	}

	if user.Email().Value() != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email().Value())
	}

	// 年龄计算可能因当前时间而异，这里只验证年龄是合理的
	age := user.Age()
	if age < 30 || age > 40 {
		t.Errorf("Expected age between 30-40, got %d", age)
	}

	if user.Status() != UserStatusActive {
		t.Errorf("Expected status Active, got %v", user.Status())
	}

	// 验证领域事件
	events := user.GetUncommittedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 uncommitted event, got %d", len(events))
	}

	if events[0].EventType() != "UserCreated" {
		t.Errorf("Expected event type 'UserCreated', got '%s'", events[0].EventType())
	}
}

func TestUser_ChangeEmail(t *testing.T) {
	// 创建用户
	userID, _ := NewUserID("user-123")
	user, _ := NewUser(userID, "testuser", "old@example.com", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))

	// 清空事件
	user.MarkEventsAsCommitted()

	// 修改邮箱
	err := user.ChangeEmail("new@example.com")
	if err != nil {
		t.Fatalf("Failed to change email: %v", err)
	}

	// 验证邮箱已更新
	if user.Email().Value() != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got '%s'", user.Email().Value())
	}

	// 验证领域事件
	events := user.GetUncommittedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 uncommitted event, got %d", len(events))
	}

	if events[0].EventType() != "UserEmailChanged" {
		t.Errorf("Expected event type 'UserEmailChanged', got '%s'", events[0].EventType())
	}
}

func TestUser_Suspend(t *testing.T) {
	// 创建用户
	userID, _ := NewUserID("user-123")
	user, _ := NewUser(userID, "testuser", "test@example.com", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))

	// 清空事件
	user.MarkEventsAsCommitted()

	// 暂停用户
	err := user.Suspend()
	if err != nil {
		t.Fatalf("Failed to suspend user: %v", err)
	}

	// 验证状态已更新
	if user.Status() != UserStatusSuspended {
		t.Errorf("Expected status Suspended, got %v", user.Status())
	}

	// 验证领域事件
	events := user.GetUncommittedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 uncommitted event, got %d", len(events))
	}

	if events[0].EventType() != "UserSuspended" {
		t.Errorf("Expected event type 'UserSuspended', got '%s'", events[0].EventType())
	}
}

func TestEmail_Validation(t *testing.T) {
	tests := []struct {
		email     string
		shouldErr bool
	}{
		{"valid@example.com", false},
		{"test.user@domain.co.uk", false},
		{"invalid-email", true},
		{"@example.com", true},
		{"test@", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			_, err := NewEmail(tt.email)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for email '%s', but got none", tt.email)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for email '%s': %v", tt.email, err)
			}
		})
	}
}

func TestBirthDate_Validation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		birthDate time.Time
		shouldErr bool
	}{
		{"Valid birth date", now.AddDate(-25, 0, 0), false},
		{"Future birth date", now.AddDate(1, 0, 0), true},
		{"Too old", now.AddDate(-200, 0, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBirthDate(tt.birthDate)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for birth date %v, but got none", tt.birthDate)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for birth date %v: %v", tt.birthDate, err)
			}
		})
	}
}
