package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/NSObjects/go-template/internal/shared/event"
	"github.com/NSObjects/go-template/internal/shared/kernel"
)

// UserID 用户ID值对象
type UserID struct {
	value string
}

// NewUserID 创建用户ID
func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}
	return UserID{value: value}, nil
}

// Value 获取ID值
func (id UserID) Value() string {
	return id.value
}

// Equals 比较ID是否相等
func (id UserID) Equals(other kernel.EntityID) bool {
	if otherUserID, ok := other.(UserID); ok {
		return id.value == otherUserID.value
	}
	return false
}

// Email 邮箱值对象
type Email struct {
	value string
}

// NewEmail 创建邮箱值对象
func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	// 简单的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: strings.ToLower(value)}, nil
}

// Value 获取邮箱值
func (e Email) Value() string {
	return e.value
}

// Equals 比较邮箱是否相等
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// BirthDate 生日值对象
type BirthDate struct {
	value time.Time
}

// NewBirthDate 创建生日值对象
func NewBirthDate(value time.Time) (BirthDate, error) {
	now := time.Now()
	if value.After(now) {
		return BirthDate{}, errors.New("birth date cannot be in the future")
	}

	// 检查年龄是否合理（0-150岁）
	age := now.Year() - value.Year()
	if age < 0 || age > 150 {
		return BirthDate{}, errors.New("invalid birth date")
	}

	return BirthDate{value: value}, nil
}

// Value 获取生日值
func (bd BirthDate) Value() time.Time {
	return bd.value
}

// Age 计算年龄
func (bd BirthDate) Age() int {
	now := time.Now()
	age := now.Year() - bd.value.Year()

	// 如果今年还没到生日，年龄减1
	if now.YearDay() < bd.value.YearDay() {
		age--
	}

	return age
}

// Equals 比较生日是否相等
func (bd BirthDate) Equals(other BirthDate) bool {
	return bd.value.Equal(other.value)
}

// User 用户聚合根
type User struct {
	kernel.BaseEntity
	kernel.BaseAggregateRoot

	username  string
	email     Email
	birthDate BirthDate
	status    UserStatus
}

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusInactive
	UserStatusSuspended
)

// NewUser 创建新用户
func NewUser(id UserID, username, emailStr string, birthDate time.Time) (*User, error) {
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	birth, err := NewBirthDate(birthDate)
	if err != nil {
		return nil, err
	}

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user := &User{
		BaseEntity: kernel.NewBaseEntity(id),
		username:   username,
		email:      email,
		birthDate:  birth,
		status:     UserStatusActive,
	}

	// 发布用户创建事件
	userCreated := &UserCreated{
		BaseEvent: *event.NewBaseEvent("UserCreated", id.Value()),
		UserID:    id.Value(),
		Username:  username,
		Email:     email.Value(),
	}
	user.AddDomainEvent(userCreated)

	return user, nil
}

// Username 获取用户名
func (u *User) Username() string {
	return u.username
}

// Email 获取邮箱
func (u *User) Email() Email {
	return u.email
}

// BirthDate 获取生日
func (u *User) BirthDate() BirthDate {
	return u.birthDate
}

// Age 获取年龄
func (u *User) Age() int {
	return u.birthDate.Age()
}

// Status 获取状态
func (u *User) Status() UserStatus {
	return u.status
}

// ChangeEmail 修改邮箱
func (u *User) ChangeEmail(newEmailStr string) error {
	if u.status != UserStatusActive {
		return errors.New("cannot change email for inactive user")
	}

	newEmail, err := NewEmail(newEmailStr)
	if err != nil {
		return err
	}

	if u.email.Equals(newEmail) {
		return errors.New("new email is the same as current email")
	}

	oldEmail := u.email.Value()
	u.email = newEmail
	u.MarkUpdated()

	// 发布邮箱变更事件
	emailChanged := &UserEmailChanged{
		BaseEvent: *event.NewBaseEvent("UserEmailChanged", u.ID().Value()),
		UserID:    u.ID().Value(),
		OldEmail:  oldEmail,
		NewEmail:  newEmail.Value(),
	}
	u.AddDomainEvent(emailChanged)

	return nil
}

// ChangeUsername 修改用户名
func (u *User) ChangeUsername(newUsername string) error {
	if u.status != UserStatusActive {
		return errors.New("cannot change username for inactive user")
	}

	if newUsername == "" {
		return errors.New("username cannot be empty")
	}

	if u.username == newUsername {
		return errors.New("new username is the same as current username")
	}

	u.username = newUsername
	u.MarkUpdated()

	return nil
}

// Suspend 暂停用户
func (u *User) Suspend() error {
	if u.status == UserStatusSuspended {
		return errors.New("user is already suspended")
	}

	u.status = UserStatusSuspended
	u.MarkUpdated()

	// 发布用户暂停事件
	userSuspended := &UserSuspended{
		BaseEvent: *event.NewBaseEvent("UserSuspended", u.ID().Value()),
		UserID:    u.ID().Value(),
		Username:  u.username,
	}
	u.AddDomainEvent(userSuspended)

	return nil
}

// Activate 激活用户
func (u *User) Activate() error {
	if u.status == UserStatusActive {
		return errors.New("user is already active")
	}

	u.status = UserStatusActive
	u.MarkUpdated()

	return nil
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.status == UserStatusActive
}
