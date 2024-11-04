package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID                int64     `bun:",pk,autoincrement"`
	Email             string    `bun:",unique,notnull"`
	Name              string    `bun:",nullzero"`
	ProfilePictureURL string    `bun:",nullzero"`
	RegistrationTime  time.Time `bun:"default:current_timestamp,notnull"`
	LastLogin         time.Time `bun:",nullzero"`
}
