package domain

type RankingParameter struct {
	ID     uint
	Alpha  float64 `json:"alpha"`  // 点赞权重
	Beta   float64 `json:"beta"`   // 收藏权重
	Gamma  float64 `json:"gamma"`  // 阅读权重
	Lambda float64 `json:"lambda"` // 时间权重
}
