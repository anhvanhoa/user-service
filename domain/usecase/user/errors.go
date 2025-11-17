package user

import (
	"github.com/anhvanhoa/service-core/domain/oops"
)

var (
	ErrGetUsers            = oops.New("Không thể lấy danh sách người dùng")
	ErrUpdateUser          = oops.New("Không thể cập nhật người dùng")
	ErrCreateUser          = oops.New("Không thể tạo người dùng")
	ErrUserAlreadyExists   = oops.New("Người dùng đã tồn tại")
	ErrEmailAlreadyExists  = oops.New("Email đã tồn tại")
	ErrPhoneAlreadyExists  = oops.New("Số điện thoại đã tồn tại")
	ErrUserAlreadyLocked   = oops.New("Người dùng đã bị khóa")
	ErrUserAlreadyUnlocked = oops.New("Người dùng đã được mở khóa")
)
