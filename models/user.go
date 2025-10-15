package models

import "github.com/google/uuid"

type UserDao struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
