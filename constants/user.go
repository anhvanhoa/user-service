package constants

import "user-service/domain/entity"

var MapStatusColumn = map[entity.UserStatus]string{
	entity.UserStatusActive:   "created_at",
	entity.UserStatusInactive: "created_at",
	entity.UserStatusLocked:   "locked_at",
}
