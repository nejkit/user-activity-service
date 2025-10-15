package models

import (
	"github.com/google/uuid"
	"time"
)

type UserActivityDao struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	ActionDate time.Time `db:"action_date"`
	Action     string    `db:"action"`
	Metadata   []byte    `db:"metadata"`
}

type EventsCountToUserDao struct {
	UserID      uuid.UUID `db:"user_id"`
	EventsCount uint64    `db:"events_count"`
}

type ActivityPeriodDao struct {
	ID       uuid.UUID `db:"id"`
	FromDate time.Time `db:"from_date"`
	ToDate   time.Time `db:"to_date"`
}

type UserActivityHistoryDao struct {
	PeriodID     uuid.UUID `db:"period_id"`
	UserID       uuid.UUID `db:"user_id"`
	ActionsCount uint64    `db:"actions_count"`
}

type ActivityHistoryDao struct {
	UserActivityHistoryDao `db:"activity"`
	ActivityPeriodDao      `db:"period"`
}
