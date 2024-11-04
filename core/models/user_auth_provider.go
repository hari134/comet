package models

import "github.com/uptrace/bun"

type UserAuthProvider struct {
	bun.BaseModel `bun:"table:user_auth_providers"`

	ID             int64 `bun:",pk,autoincrement"`
	UserID         int64 `bun:"notnull"` // Foreign key to users table
	ProviderID     int64 `bun:"notnull"` // Foreign key to auth_providers table
	ProviderAuthID int64 `bun:"notnull"` // ID from the provider-specific table
}
