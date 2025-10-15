package services

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"time"
	"user-activity-service/dto"
	"user-activity-service/models"
)

type ActivityService struct {
	userRepository            userRepository
	activityRepository        activityRepository
	activityHistoryRepository activityHistoryRepository
}

func NewActivityService(userRepository userRepository, activityRepository activityRepository, activityHistoryRepository activityHistoryRepository) *ActivityService {
	return &ActivityService{userRepository: userRepository, activityRepository: activityRepository, activityHistoryRepository: activityHistoryRepository}
}

func (a *ActivityService) SaveUserActivity(ctx context.Context, activityDto *dto.SaveUserActivityRequestDto) error {
	userID := uuid.MustParse(activityDto.UserID)

	if _, err := a.userRepository.GetByID(ctx, userID); err != nil {
		return err
	}

	if activityDto.Metadata == nil {
		activityDto.Metadata = make(map[string]interface{})
	}

	metadataBytes, _ := json.Marshal(activityDto.Metadata)

	activityInfo := &models.UserActivityDao{
		ID:         uuid.New(),
		UserID:     userID,
		ActionDate: time.Now().UTC(),
		Action:     activityDto.Action,
		Metadata:   metadataBytes,
	}

	return a.activityRepository.Save(ctx, activityInfo)
}

func (a *ActivityService) GetUserActivities(ctx context.Context, query *dto.GetUserActivitiesQuery) (*dto.GetUserActivityResponseDto, error) {
	userID := uuid.MustParse(query.UserID)

	if _, err := a.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	activities, err := a.activityRepository.GetByUserID(ctx, userID, query.FromDate, query.ToDate)

	if err != nil {
		return nil, err
	}

	var mappedResult = new(dto.GetUserActivityResponseDto)
	mappedResult.Activities = make([]dto.UserActivityEventDto, len(activities))

	for i, activity := range activities {
		parsedMetadata := make(map[string]interface{})

		if err = json.Unmarshal(activity.Metadata, &parsedMetadata); err != nil {
			return nil, err
		}

		mappedResult.Activities[i] = dto.UserActivityEventDto{
			ActionDate: activity.ActionDate,
			Action:     activity.Action,
			Metadata:   parsedMetadata,
		}
	}

	return mappedResult, nil
}

func (a *ActivityService) GetUsersActivities(ctx context.Context, query *dto.GetActivitiesQuery) (*dto.GetUsersActivityResponseDto, error) {
	activities, err := a.activityRepository.GetAll(ctx, query.FromDate, query.ToDate)

	if err != nil {
		return nil, err
	}

	var mappedResult = new(dto.GetUsersActivityResponseDto)
	mappedResult.UserActivities = make(map[uuid.UUID][]dto.UserActivityEventDto)

	for _, activity := range activities {
		currentEvents, ok := mappedResult.UserActivities[activity.UserID]

		if !ok {
			currentEvents = make([]dto.UserActivityEventDto, 0)
			mappedResult.UserActivities[activity.UserID] = currentEvents
		}

		parsedMetadata := make(map[string]interface{})

		if err = json.Unmarshal(activity.Metadata, &parsedMetadata); err != nil {
			return nil, err
		}

		currentEvents = append(currentEvents, dto.UserActivityEventDto{
			ActionDate: activity.ActionDate,
			Action:     activity.Action,
			Metadata:   parsedMetadata,
		})

		mappedResult.UserActivities[activity.UserID] = currentEvents
	}

	return mappedResult, nil
}

func (a *ActivityService) GetUserActivityHistories(ctx context.Context, query *dto.GetUserActivitiesQuery) (interface{}, error) {
	userID := uuid.MustParse(query.UserID)

	if _, err := a.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	activityHistories, err := a.activityHistoryRepository.GetByUserID(ctx, userID, query.FromDate, query.ToDate)

	if err != nil {
		return nil, err
	}

	var mappedResult = new(dto.GetUserActivityStatisticResponseDto)
	mappedResult.Statistics = make([]dto.UserActivityStatisticDto, len(activityHistories))

	for i, activity := range activityHistories {
		mappedResult.Statistics[i] = dto.UserActivityStatisticDto{
			FromDate:     activity.FromDate,
			ToDate:       activity.ToDate,
			ActionsCount: activity.ActionsCount,
		}
	}

	return mappedResult, nil
}

func (a *ActivityService) GetUsersActivityHistories(ctx context.Context, query *dto.GetActivitiesQuery) (interface{}, error) {
	activityHistories, err := a.activityHistoryRepository.GetAll(ctx, query.FromDate, query.ToDate)

	if err != nil {
		return nil, err
	}

	var mappedResult = new(dto.GetUsersActivityStatisticResponseDto)
	mappedResult.UserActivities = make(map[uuid.UUID][]dto.UserActivityStatisticDto)

	for _, activity := range activityHistories {
		currentEvents, ok := mappedResult.UserActivities[activity.UserID]

		if !ok {
			currentEvents = make([]dto.UserActivityStatisticDto, 0)
			mappedResult.UserActivities[activity.UserID] = currentEvents
		}

		currentEvents = append(currentEvents, dto.UserActivityStatisticDto{
			FromDate:     activity.FromDate,
			ToDate:       activity.ToDate,
			ActionsCount: activity.ActionsCount,
		})

		mappedResult.UserActivities[activity.UserID] = currentEvents
	}

	return mappedResult, nil
}
