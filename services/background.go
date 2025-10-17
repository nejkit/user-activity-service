package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
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
				logger.Infoln("process next period without pause")
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

	logger.
		WithField("lastPeriodDate", lastPeriodDate).
		Infoln("calculate activity statistic")

	if err != nil {
		logger.WithError(err).Errorln("failed get last period date")
		return false
	}

	fromDate := *lastPeriodDate

	if lastPeriodDate.Add(b.cfg.EventsPeriodDuration).After(nowDate) {
		logger.Infoln("activity period not be ended for calculate, cancel processing")
		return false
	}

	if err = b.calculateEventsCount(ctx, fromDate, fromDate.Add(b.cfg.EventsPeriodDuration)); err != nil {
		logger.WithError(err).Errorln("failed process new activity period")
		return false
	}

	logger.Infoln("activity period processed successfully")

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

	eventsCount, err := b.activityRepository.GetEventsCount(ctx, fromDate, toDate)

	if err != nil {
		return err
	}

	if err = b.activityPeriodsRepository.Save(ctx, newPeriod); err != nil {
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

	if len(periodActivities) == 0 {
		logger.Infoln("period has not user activity, save empty period")
		return nil
	}

	return b.activityHistoryRepository.SaveAll(ctx, periodActivities)
}
