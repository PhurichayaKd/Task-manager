package domain

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64          `db:"id"`
	Email        string         `db:"email"`
	Username     sql.NullString `db:"username"`
	PasswordHash sql.NullString `db:"password_hash"`
	Role         string         `db:"role"`

	Name       sql.NullString `db:"name"`
	Provider   sql.NullString `db:"provider"`
	ProviderID sql.NullString `db:"provider_id"`
	AvatarURL  sql.NullString `db:"avatar_url"`

	CreatedAt time.Time `db:"created_at"`
}
