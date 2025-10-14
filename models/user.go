package models

import "github.com/google/uuid"

type UserDao struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
