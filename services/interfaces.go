package services

import (
	"context"
	"github.com/google/uuid"
	"time"
	"user-activity-service/models"
)

type userRepository interface {
	Save(ctx context.Context, dao *models.UserDao) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.UserDao, error)
	DeleteByID(ctx context.Context, userID uuid.UUID) error
}

type activityRepository interface {
	Save(ctx context.Context, dao *models.UserActivityDao) error
	GetAll(ctx context.Context, fromDate, toDate *time.Time) ([]models.UserActivityDao, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate *time.Time) ([]models.UserActivityDao, error)
	GetEventsCount(ctx context.Context, fromDate, toDate time.Time) ([]models.EventsCountToUserDao, error)
}

type activityHistoryRepository interface {
	SaveAll(ctx context.Context, activityHistories []models.UserActivityHistoryDao) error
	GetByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate *time.Time) ([]models.ActivityHistoryDao, error)
	GetAll(ctx context.Context, fromDate, toDate *time.Time) ([]models.ActivityHistoryDao, error)
}

type activityPeriodsRepository interface {
	Save(ctx context.Context, dao *models.ActivityPeriodDao) error
	GetLastActivityPeriod(ctx context.Context) (*models.ActivityPeriodDao, error)
}
