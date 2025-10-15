package dto

import (
	"github.com/google/uuid"
	"time"
)

type RegisterUserResponseDto struct {
	UserID uuid.UUID `json:"userId"`
}

type UserActivityEventDto struct {
	ActionDate time.Time              `json:"actionDate"`
	Action     string                 `json:"action"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type GetUserActivityResponseDto struct {
	Activities []UserActivityEventDto `json:"activities"`
}

type GetUsersActivityResponseDto struct {
	UserActivities map[uuid.UUID][]UserActivityEventDto `json:"userActivities"`
}

type UserActivityStatisticDto struct {
	FromDate     time.Time `json:"fromDate"`
	ToDate       time.Time `json:"toDate"`
	ActionsCount uint64    `json:"actionsCount"`
}

type GetUserActivityStatisticResponseDto struct {
	Statistics []UserActivityStatisticDto `json:"statistics"`
}

type GetUsersActivityStatisticResponseDto struct {
	UserActivities map[uuid.UUID][]UserActivityStatisticDto `json:"userActivities"`
}
