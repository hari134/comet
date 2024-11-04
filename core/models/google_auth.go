package models

import "github.com/uptrace/bun"

type GoogleAuth struct {
	bun.BaseModel `bun:"table:google_auth"`

	ID       int64  `bun:",pk,autoincrement"`
	GoogleID string `bun:",unique,notnull"`
}
