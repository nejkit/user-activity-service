package dto

import "time"

type RegisterUserRequestDto struct {
	Name string `json:"name" validate:"required"`
}

type SaveUserActivityRequestDto struct {
	UserID   string                 `json:"userId" validate:"required,uuid"`
	Action   string                 `json:"action" validate:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

type GetUserActivitiesQuery struct {
	UserID string `query:"userId" validate:"required,uuid"`
	GetActivitiesQuery
}

type GetActivitiesQuery struct {
	FromDate *time.Time `query:"fromDate,omitempty"`
	ToDate   *time.Time `query:"toDate,omitempty"`
}
