package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"time"
	"user-activity-service/models"
	"user-activity-service/services"
)

type PgUserStorage struct {
	client *sqlx.DB
}

func NewPgUserStorage(db *sqlx.DB) *PgUserStorage {
	return &PgUserStorage{client: db}
}

func (p *PgUserStorage) Save(ctx context.Context, dao *models.UserDao) error {
	cmd, args := sqlbuilder.InsertInto(userTableName).
		Cols("id", "name").
		Values(dao.ID, dao.Name).
		Build()

	if _, err := p.client.ExecContext(ctx, cmd, args...); err != nil {
		return errors.Join(services.ErrorInternalError, err)
	}

	return nil
}

func (p *PgUserStorage) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserDao, error) {
	queryBuilder := sqlbuilder.Select("*").
		From(userTableName)

	queryBuilder.Where(queryBuilder.Equal("id", userID))

	query, args := queryBuilder.Build()

	var dao = new(models.UserDao)

	err := p.client.GetContext(ctx, dao, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, services.ErrorUserNotFound
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return dao, nil
}

func (p *PgUserStorage) DeleteByID(ctx context.Context, userID uuid.UUID) error {
	builder := sqlbuilder.DeleteFrom(userTableName)

	builder.Where(builder.Equal("id", userID))

	query, args := builder.Build()

	if _, err := p.client.ExecContext(ctx, query, args...); err != nil {
		return errors.Join(services.ErrorInternalError, err)
	}

	return nil
}

type PgUserActivityStorage struct {
	client *sqlx.DB
}

func NewPgUserActivityStorage(client *sqlx.DB) *PgUserActivityStorage {
	return &PgUserActivityStorage{client: client}
}

func (p *PgUserActivityStorage) Save(ctx context.Context, dao *models.UserActivityDao) error {
	cmd, args := sqlbuilder.InsertInto(userActivityTableName).
		Cols("id", "user_id", "action_date", "action", "metadata").
		Values(dao.ID, dao.UserID, dao.ActionDate, dao.Action, dao.Metadata).
		Build()

	if _, err := p.client.ExecContext(ctx, cmd, args...); err != nil {
		return errors.Join(services.ErrorInternalError, err)
	}

	return nil
}

func (p *PgUserActivityStorage) GetAll(ctx context.Context, fromDate, toDate *time.Time) ([]models.UserActivityDao, error) {
	queryBuilder := sqlbuilder.Select("*").
		From(userActivityTableName)

	if fromDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.GTE("action_date", *fromDate))
	}

	if toDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.LTE("action_date", *toDate))
	}

	query, args := queryBuilder.Build()

	var activities []models.UserActivityDao

	err := p.client.SelectContext(ctx, &activities, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return make([]models.UserActivityDao, 0), nil
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return activities, nil
}

func (p *PgUserActivityStorage) GetByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate *time.Time) ([]models.UserActivityDao, error) {
	queryBuilder := sqlbuilder.Select("*").
		From(userActivityTableName)

	queryBuilder = queryBuilder.Where(queryBuilder.Equal("user_id", userID))

	if fromDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.GTE("action_date", *fromDate))
	}

	if toDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.LTE("action_date", *toDate))
	}

	query, args := queryBuilder.Build()

	var activities []models.UserActivityDao

	err := p.client.SelectContext(ctx, &activities, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return make([]models.UserActivityDao, 0), nil
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return activities, nil
}

func (p *PgUserActivityStorage) GetEventsCount(ctx context.Context, fromDate, toDate time.Time) ([]models.EventsCountToUserDao, error) {
	queryBuilder := sqlbuilder.Select("user_id", "count(*) as \"events_count\"").
		From(userActivityTableName)

	query, args := queryBuilder.Where(
		queryBuilder.GTE("action_date", fromDate),
		queryBuilder.LTE("action_date", toDate),
	).
		GroupBy("user_id").
		Build()

	var result []models.EventsCountToUserDao

	err := p.client.SelectContext(ctx, &result, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return make([]models.EventsCountToUserDao, 0), nil
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return result, nil
}

type PgUserActivityHistoryStorage struct {
	client *sqlx.DB
}

func NewPgUserActivityHistoryStorage(client *sqlx.DB) *PgUserActivityHistoryStorage {
	return &PgUserActivityHistoryStorage{client: client}
}

func (p *PgUserActivityHistoryStorage) SaveAll(ctx context.Context, dao []models.UserActivityHistoryDao) error {
	builder := sqlbuilder.InsertInto(userActivityHistoryTableName).
		Cols("period_id", "user_id", "actions_count")

	for _, event := range dao {
		builder = builder.Values(event.PeriodID, event.UserID, event.ActionsCount)
	}

	cmd, args := builder.Build()

	if _, err := p.client.ExecContext(ctx, cmd, args...); err != nil {
		return errors.Join(services.ErrorInternalError, err)
	}

	return nil
}

func (p *PgUserActivityHistoryStorage) GetByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate *time.Time) ([]models.ActivityHistoryDao, error) {
	queryBuilder := sqlbuilder.Select("period.id, period.from_date, period.to_date, activity.period_id, activity.user_id, activity.actions_count")

	queryBuilder = queryBuilder.
		From(queryBuilder.As(userActivityTableName, "activity")).
		JoinWithOption(sqlbuilder.InnerJoin, queryBuilder.As(activityPeriodsTableName, "period"), "(period.id = activity.period_id)").
		Where(queryBuilder.Equal("activity.user_id", userID))

	if fromDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.GTE("period.from_date", *fromDate))
	}

	if toDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.LTE("period.to_date", *toDate))
	}

	query, args := queryBuilder.Build()

	var activities []models.ActivityHistoryDao

	err := p.client.SelectContext(ctx, &activities, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return make([]models.ActivityHistoryDao, 0), nil
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return activities, nil
}

func (p *PgUserActivityHistoryStorage) GetAll(ctx context.Context, fromDate, toDate *time.Time) ([]models.ActivityHistoryDao, error) {
	queryBuilder := sqlbuilder.Select("period.id, period.from_date, period.to_date, activity.period_id, activity.user_id, activity.actions_count")

	queryBuilder = queryBuilder.
		From(queryBuilder.As(userActivityTableName, "activity")).
		JoinWithOption(sqlbuilder.InnerJoin, queryBuilder.As(activityPeriodsTableName, "period"), "(period.id = activity.period_id)")

	if fromDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.GTE("period.from_date", *fromDate))
	}

	if toDate != nil {
		queryBuilder = queryBuilder.Where(queryBuilder.LTE("period.to_date", *toDate))
	}

	query, args := queryBuilder.Build()

	var activities []models.ActivityHistoryDao

	err := p.client.SelectContext(ctx, &activities, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return make([]models.ActivityHistoryDao, 0), nil
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return activities, nil
}

type PgActivityPeriodsStorage struct {
	client *sqlx.DB
}

func NewPgActivityPeriodsStorage(client *sqlx.DB) *PgActivityPeriodsStorage {
	return &PgActivityPeriodsStorage{client: client}
}

func (p *PgActivityPeriodsStorage) Save(ctx context.Context, dao *models.ActivityPeriodDao) error {
	cmd, args := sqlbuilder.InsertInto(activityPeriodsTableName).
		Cols("id", "from_date", "to_date").
		Values(dao.ID, dao.FromDate, dao.ToDate).
		Build()

	if _, err := p.client.ExecContext(ctx, cmd, args...); err != nil {
		return errors.Join(services.ErrorInternalError, err)
	}

	return nil
}

func (p *PgActivityPeriodsStorage) GetLastActivityPeriod(ctx context.Context) (*models.ActivityPeriodDao, error) {
	builder := sqlbuilder.Select("*").
		From(activityPeriodsTableName)

	query, args := builder.OrderByDesc("to_date").
		Limit(1).
		Build()

	var dao = new(models.ActivityPeriodDao)

	err := p.client.SelectContext(ctx, &dao, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, services.ErrorActivityPeriodNotFound
	}

	if err != nil {
		return nil, errors.Join(services.ErrorInternalError, err)
	}

	return dao, nil
}
