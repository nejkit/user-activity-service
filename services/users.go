package services

import (
	"context"
	"github.com/google/uuid"
	"user-activity-service/models"
)

type UserService struct {
	userRepository userRepository
}

func NewUserService(userRepository userRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) Register(ctx context.Context, name string) (*uuid.UUID, error) {
	user := &models.UserDao{
		ID:   uuid.New(),
		Name: name,
	}

	if err := u.userRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (u *UserService) GetUserName(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := u.userRepository.GetByID(ctx, userID)

	if err != nil {
		return "", err
	}

	return user.Name, nil
}

func (u *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return u.userRepository.DeleteByID(ctx, userID)
}
