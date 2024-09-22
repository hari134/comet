package models

import "github.com/uptrace/bun"

type AuthProvider struct {
    bun.BaseModel `bun:"table:auth_providers"`

    ID   int64  `bun:",pk,autoincrement"`
    Name string `bun:",unique,notnull"`
}
