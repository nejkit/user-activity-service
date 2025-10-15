package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"user-activity-service/dto"
	"user-activity-service/services"
)

type UserHandler struct {
	userService *services.UserService
	Wg          *sync.WaitGroup
	ApiStopped  bool
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService, Wg: &sync.WaitGroup{}, ApiStopped: false}
}

func (u *UserHandler) HandleRegisterUser(ctx *gin.Context) {
	if u.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	u.Wg.Add(1)
	defer u.Wg.Done()

	var request = new(dto.RegisterUserRequestDto)

	if err := ctx.ShouldBindJSON(request); err != nil {
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
	}

	userID, err := u.userService.Register(ctx, request.Name)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, dto.RegisterUserResponseDto{UserID: *userID})
}

func (u *UserHandler) HandleRemoveUser(ctx *gin.Context) {
	if u.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	u.Wg.Add(1)
	defer u.Wg.Done()

	userIDParam := ctx.Param("id")

	userID, err := uuid.Parse(userIDParam)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: "malformed id",
		})
		return
	}

	if err = u.userService.Delete(ctx, userID); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func (u *UserHandler) HandleGetUserName(ctx *gin.Context) {
	if u.ApiStopped {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	u.Wg.Add(1)
	defer u.Wg.Done()

	userIDParam := ctx.Param("id")

	userID, err := uuid.Parse(userIDParam)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: "malformed id",
		})
		return
	}

	username, err := u.userService.GetUserName(ctx, userID)

	if errors.Is(err, services.ErrorUserNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponseDto{
			Code:    dto.ErrorCodeAccountNotFound,
			Message: "user not found",
		})
		return
	}

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"username": username})
}
