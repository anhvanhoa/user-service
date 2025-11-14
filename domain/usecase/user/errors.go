package user

import (
	"github.com/anhvanhoa/service-core/domain/oops"
)

var (
	ErrGetUsers = oops.New("Không thể lấy danh sách người dùng")
)
