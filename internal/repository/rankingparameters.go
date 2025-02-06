package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type RankingParameterRepository interface {
	Insert(ctx context.Context, rankingParameter domain.RankingParameter) (uint, error) // Note:为了保证可以看到历史参数的内容，只以追加的形式插入，不进行更新
	FindLastParameter(ctx context.Context) (domain.RankingParameter, error)             // Note: 查找最后一次插入的参数
}
type rankingParameterRepository struct {
	dao dao.RankingParameterDAO
	l   *zap.Logger
}

func NewRankingParameterRepository(dao dao.RankingParameterDAO, l *zap.Logger) RankingParameterRepository {
	return &rankingParameterRepository{
		dao: dao,
		l:   l,
	}
}

func (r *rankingParameterRepository) Insert(ctx context.Context, rankingParameter domain.RankingParameter) (uint, error) {
	rankingParameterID, err := r.dao.Insert(ctx, toDaoRankingParameter(rankingParameter))
	if err != nil {
		r.l.Error("插入RankingParameter失败", zap.Error(err))
		return 0, err
	}
	return rankingParameterID, nil
}
func (r *rankingParameterRepository) FindLastParameter(ctx context.Context) (domain.RankingParameter, error) {
	rankingParameterVal, err := r.dao.FindLastParameter(ctx)
	if err != nil {
		r.l.Error("查找RankingParameter失败", zap.Error(err))
		return domain.RankingParameter{}, err
	}
	return toDomainRankingParameter(rankingParameterVal), nil
}
func toDomainRankingParameter(rankingParameter dao.RankingParameter) domain.RankingParameter {
	return domain.RankingParameter{
		ID:     rankingParameter.ID,
		Alpha:  rankingParameter.Alpha,
		Beta:   rankingParameter.Beta,
		Gamma:  rankingParameter.Gamma,
		Lambda: rankingParameter.Lambda,
	}
}
func toDaoRankingParameter(rankingParameter domain.RankingParameter) dao.RankingParameter {
	return dao.RankingParameter{
		ID:     rankingParameter.ID,
		Alpha:  rankingParameter.Alpha,
		Beta:   rankingParameter.Beta,
		Gamma:  rankingParameter.Gamma,
		Lambda: rankingParameter.Lambda,
	}
}
