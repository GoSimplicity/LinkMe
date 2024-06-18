package api

type CreatePlateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeletePlateReq struct {
	ID int64 `uri:"id"`
}

type UpdatePlateReq struct {
	ID int64 `uri:"id"`
}
