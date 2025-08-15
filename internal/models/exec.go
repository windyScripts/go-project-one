package models

import "database/sql"

type Exec struct {
	ID                  int            `json:"id,omitempty" db:"id,omitempty"`
	FirstName           string         `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName            string         `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email               string         `json:"email,omitempty" db:"email,omitempty"`
	Username            string         `json:"username,omitempty" db:"username,omitempty"`
	Password            string         `json:"password,omitempty" db:"password,omitempty"`
	PasswordChangedAt   sql.NullString `json:"password_changed_at,omitempty" db:"password_changed_at,omitempty"`
	UserCreatedAt       sql.NullString `json:"user_created_at,omitempty" db:"user_created_at,omitempty"`
	PasswordResetCode   sql.NullString `json:"password_reset_code,omitempty" db:"password_reset_code,omitempty"`
	PasswordCodeExpires sql.NullString `json:"password_code_expires,omitempty" db:"password_code_expires,omitempty"`
	InactiveStatus      bool           `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
	Role                string         `json:"role,omitempty" db:"role,omitempty"`
}

/*
bcrypt, argon2 and pbkdf2 are hashing algos.
bcrypt is well established. efficient enough. widely supported.
pbkdf2 password based key derivative function2 secure, but weak to side channel attacks etc.
widely supported, standardized.
argon2 argon2d, argon2i, argon2id. Efficient in time and parallelism. New, not as supported as bcrypt. Most secure, flexible,
followed by bcrypt and then pbkdf2
*/
