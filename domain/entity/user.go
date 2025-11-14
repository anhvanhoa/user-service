package entity

import (
	"time"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusDeleted  UserStatus = "deleted"
)

type User struct {
	tableName  struct{}   `pg:"users,alias:u"`
	ID         string     `pg:"id,pk"`
	Email      string     `pg:"email,unique"`
	Phone      string     `pg:"phone,unique"`
	Password   string     `pg:"password"`
	FullName   string     `pg:"full_name"`
	Avatar     string     `pg:"avatar"`
	Bio        string     `pg:"bio"`
	Address    string     `pg:"address"`
	CodeVerify string     `pg:"code_verify"`
	Veryfied   *time.Time `pg:"veryfied"`
	CreatedBy  string     `pg:"created_by"`
	Status     UserStatus `pg:"status"`
	Birthday   *time.Time `pg:"birthday"`
	DeletedAt  *time.Time `pg:"deleted_at"`
	CreatedAt  time.Time  `pg:"created_at"`
	UpdatedAt  *time.Time `pg:"updated_at"`
}

type UserInfor struct {
	ID       string
	Email    string
	Phone    string
	FullName string
	Avatar   string
	Bio      string
	Address  string
	Birthday *time.Time
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetNameTable() any {
	return u.tableName
}

func (u *User) GetInfor() UserInfor {
	return UserInfor{
		ID:       u.ID,
		Email:    u.Email,
		Phone:    u.Phone,
		FullName: u.FullName,
		Avatar:   u.Avatar,
		Bio:      u.Bio,
		Address:  u.Address,
		Birthday: u.Birthday,
	}
}

type FilterUser struct {
	Status   *UserStatus
	FromDate *time.Time
	ToDate   *time.Time
}
