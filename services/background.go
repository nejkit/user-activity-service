package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
	"user-activity-service/config"
	"user-activity-service/models"
)

type BackgroundActivityService struct {
	activityPeriodsRepository activityPeriodsRepository
	activityHistoryRepository activityHistoryRepository
	activityRepository        activityRepository
	cfg                       *config.BackgroundWorkerConfig
	Wg                        *sync.WaitGroup
}

func NewBackgroundActivityService(activityPeriodsRepository activityPeriodsRepository, activityHistoryRepository activityHistoryRepository, activityRepository activityRepository, cfg *config.BackgroundWorkerConfig) *BackgroundActivityService {
	return &BackgroundActivityService{activityPeriodsRepository: activityPeriodsRepository, activityHistoryRepository: activityHistoryRepository, activityRepository: activityRepository, cfg: cfg, Wg: &sync.WaitGroup{}}
}

func (b *BackgroundActivityService) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			b.Wg.Wait()
			return

		default:
			if withRetry := b.process(ctx); withRetry {
				continue
			}

			time.Sleep(b.cfg.Interval)
		}
	}
}

func (b *BackgroundActivityService) process(ctx context.Context) bool {
	b.Wg.Add(1)
	defer b.Wg.Done()

	nowDate := time.Now().UTC()
	lastPeriodDate, err := b.initLastCalculatedPeriodDate(ctx)

	if err != nil {
		//TODO: logging
		return false
	}

	fromDate := *lastPeriodDate

	if lastPeriodDate.Add(b.cfg.EventsPeriodDuration).After(nowDate) {
		//TODO: logging
		return false
	}

	if err = b.calculateEventsCount(ctx, fromDate, fromDate.Add(b.cfg.EventsPeriodDuration)); err != nil {
		return false
	}

	return lastPeriodDate.Add(b.cfg.EventsPeriodDuration).Before(nowDate)
}

func (b *BackgroundActivityService) initLastCalculatedPeriodDate(ctx context.Context) (*time.Time, error) {
	lastPeriod, err := b.activityPeriodsRepository.GetLastActivityPeriod(ctx)

	if errors.Is(err, ErrorActivityPeriodNotFound) {
		lastPeriod = &models.ActivityPeriodDao{
			ID:       uuid.New(),
			FromDate: time.UnixMilli(0),
			ToDate:   time.Now().UTC(),
		}

		if err = b.activityPeriodsRepository.Save(ctx, lastPeriod); err != nil {
			return nil, err
		}

		return &lastPeriod.ToDate, nil
	}

	if err != nil {
		return nil, err
	}

	return &lastPeriod.ToDate, nil
}

func (b *BackgroundActivityService) calculateEventsCount(ctx context.Context, fromDate, toDate time.Time) error {
	newPeriod := &models.ActivityPeriodDao{
		ID:       uuid.New(),
		FromDate: fromDate,
		ToDate:   toDate,
	}

	if err := b.activityPeriodsRepository.Save(ctx, newPeriod); err != nil {
		return err
	}

	eventsCount, err := b.activityRepository.GetEventsCount(ctx, fromDate, toDate)

	if err != nil {
		return err
	}

	var periodActivities = make([]models.UserActivityHistoryDao, len(eventsCount))

	for i, event := range eventsCount {
		periodActivities[i] = models.UserActivityHistoryDao{
			PeriodID:     newPeriod.ID,
			UserID:       event.UserID,
			ActionsCount: event.EventsCount,
		}
	}

	return b.activityHistoryRepository.SaveAll(ctx, periodActivities)
}
