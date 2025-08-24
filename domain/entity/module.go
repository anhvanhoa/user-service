package entity

import (
	"auth-service/domain/common"
	"time"
)

type Module struct {
	tableName   struct{}      `pg:"modules,alias:m"`
	ID          string        `pg:"id,pk"`
	Name        string        `pg:"name"`
	Description string        `pg:"description"`
	Status      common.Status `pg:"status"`
	CreatedAt   time.Time     `pg:"created_at"`
	UpdatedAt   *time.Time    `pg:"updated_at"`
}

func (m *Module) NameTable() any {
	return m.tableName
}
