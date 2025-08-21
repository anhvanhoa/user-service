package entity

import (
	"time"
)

type ModuleChildRole struct {
	tableName     struct{}   `pg:"module_child_roles,alias:mcr"`
	RoleID        string     `pg:"role_id,pk"`
	ModuleChildID string     `pg:"module_child_id,pk"`
	CreatedBy     string     `pg:"created_by"`
	CreatedAt     time.Time  `pg:"created_at"`
	UpdatedAt     *time.Time `pg:"updated_at"`
}

func (mcr *ModuleChildRole) NameTable() any {
	return mcr.tableName
}
