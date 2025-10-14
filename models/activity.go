package models

import (
	"github.com/google/uuid"
	"time"
)

type UserActivityDao struct {
	Id         uuid.UUID              `db:"id"`
	UserId     uuid.UUID              `db:"user_id"`
	ActionDate time.Time              `db:"action_date"`
	Action     string                 `db:"action"`
	Metadata   map[string]interface{} `db:"metadata"`
}

type UserActivityHistoryDao struct {
	Id           uuid.UUID `db:"id"`
	UserId       uuid.UUID `db:"user_id"`
	FromDate     time.Time `db:"from_date"`
	ToDate       time.Time `db:"to_date"`
	ActionsCount int64     `db:"actions_count"`
}
