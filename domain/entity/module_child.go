package entity

import (
	"auth-service/domain/common"
	"time"
)

type ModuleChild struct {
	tableName struct{}      `pg:"module_childs,alias:mc"`
	ID        string        `pg:"id,pk"`
	ModuleID  string        `pg:"module_id"`
	Name      string        `pg:"name"`
	Path      string        `pg:"path"`
	Method    string        `pg:"method"`
	IsPrivate bool          `pg:"is_private"`
	Status    common.Status `pg:"status"`
	CreatedAt time.Time     `pg:"created_at"`
	UpdatedAt *time.Time    `pg:"updated_at"`
}

func (mc *ModuleChild) NameTable() any {
	return mc.tableName
}
