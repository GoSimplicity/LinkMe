package dao

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type RankingParameterDAO interface {
	Insert(ctx context.Context, rankingParameter RankingParameter) (uint, error) // Note:为了保证可以看到历史参数的内容，只以追加的形式插入，不进行更新
	FindLastParameter(ctx context.Context) (RankingParameter, error)             // Note: 查找最后一次插入的参数
	//Update(ctx context.Context, rankingParameter RankingParameter) error
	//GetById(ctx context.Context, id uint) (RankingParameter, error)
}
type rankingParameterDAO struct {
	l  *zap.Logger
	db *gorm.DB
}

type RankingParameter struct {
	gorm.Model
	ID     uint    `gorm:"primaryKey" json:"id"`               // 主键 ID
	Alpha  float64 `gorm:"not null;default:1.0" json:"alpha"`  // 默认值 1.0
	Beta   float64 `gorm:"not null;default:10.0" json:"beta"`  // 默认值 10.0
	Gamma  float64 `gorm:"not null;default:20.0" json:"gamma"` // 默认值 20.0
	Lambda float64 `gorm:"not null;default:1.2" json:"lambda"` // 默认值 1.2
	// CreatedAt 和 UpdatedAt 可根据需要自动填充
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func NewRankingParameterDAO(db *gorm.DB, l *zap.Logger) RankingParameterDAO {
	return &rankingParameterDAO{
		l:  l,
		db: db,
	}
}
func (r *rankingParameterDAO) Insert(ctx context.Context, rankingParameter RankingParameter) (uint, error) {
	if err := r.db.WithContext(ctx).Create(&rankingParameter).Error; err != nil {
		r.l.Error("插入RankingParameter失败", zap.Error(err))
		return 0, err
	}
	return rankingParameter.ID, nil
}
func (r *rankingParameterDAO) FindLastParameter(ctx context.Context) (RankingParameter, error) {
	var rankingParameter RankingParameter
	if err := r.db.WithContext(ctx).Last(&rankingParameter).Error; err != nil {
		r.l.Error("查找RankingParameter失败", zap.Error(err))
		return RankingParameter{}, err
	}
	return rankingParameter, nil
}
