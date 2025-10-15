package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"user-activity-service/dto"
	"user-activity-service/services"
)

type ActivityHandler struct {
	activityService *services.ActivityService
	Wg              *sync.WaitGroup
	ApiStopped      bool
}

func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService, Wg: &sync.WaitGroup{}, ApiStopped: false}
}

func (a *ActivityHandler) HandleRegisterActivity(ctx *gin.Context) {
	if a.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	a.Wg.Add(1)
	defer a.Wg.Done()

	var request = new(dto.SaveUserActivityRequestDto)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: "malformed json",
		})
		return
	}

	if err := validateRequest(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	if err := a.activityService.SaveUserActivity(ctx, request); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (a *ActivityHandler) HandleGetActivityEvents(ctx *gin.Context) {
	if a.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	a.Wg.Add(1)
	defer a.Wg.Done()

	var query = new(dto.GetActivitiesQuery)

	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	if err := validateRequest(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := a.activityService.GetUsersActivities(ctx, query)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *ActivityHandler) HandleGetUserActivityEvents(ctx *gin.Context) {
	if a.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	a.Wg.Add(1)
	defer a.Wg.Done()

	var query = new(dto.GetUserActivitiesQuery)

	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	if err := validateRequest(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := a.activityService.GetUserActivities(ctx, query)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *ActivityHandler) HandleGetActivityStatistic(ctx *gin.Context) {
	if a.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	a.Wg.Add(1)
	defer a.Wg.Done()

	var query = new(dto.GetActivitiesQuery)

	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	if err := validateRequest(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := a.activityService.GetUsersActivityHistories(ctx, query)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *ActivityHandler) HandleGetUserActivityStatistic(ctx *gin.Context) {
	if a.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	a.Wg.Add(1)
	defer a.Wg.Done()

	var query = new(dto.GetUserActivitiesQuery)

	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	if err := validateRequest(query); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := a.activityService.GetUserActivityHistories(ctx, query)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
