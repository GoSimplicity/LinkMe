package service

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
)

type ActivityService interface {
	GetRecentActivity(ctx context.Context) ([]domain.RecentActivity, error)
}

type activityService struct {
	repo repository.ActivityRepository
}

func NewActivityService(repo repository.ActivityRepository) ActivityService {
	return &activityService{
		repo: repo,
	}
}

func (a *activityService) GetRecentActivity(ctx context.Context) ([]domain.RecentActivity, error) {
	activity, err := a.repo.GetRecentActivity(ctx)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
