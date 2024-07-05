package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
)

type ActivityService interface {
	GetRecentActivity(ctx context.Context) (domain.RecentActivity, error)
}

type activityService struct {
	repo repository.ActivityRepository
}

func NewActivityService(repo repository.ActivityRepository) ActivityService {
	return &activityService{
		repo: repo,
	}
}

func (a *activityService) GetRecentActivity(ctx context.Context) (domain.RecentActivity, error) {
	activity, err := a.repo.GetRecentActivity(ctx)
	if err != nil {
		return domain.RecentActivity{}, err
	}
	return activity, nil
}
