package req

type RankingParameterReq struct {
	ID     uint    `json:"id"`
	Alpha  float64 `json:"alpha"`
	Beta   float64 `json:"beta" `
	Gamma  float64 `json:"gamma"`
	Lambda float64 `json:"lambda"`
}
