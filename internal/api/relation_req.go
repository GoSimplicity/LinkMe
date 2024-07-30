package api

type ListRelationsReq struct {
	FollowerID int64  `json:"followerId"`
	Page       int    `json:"page,omitempty"` // 当前页码
	Size       *int64 `json:"size,omitempty"` // 每页数据量
}

type GetRelationInfoReq struct {
	FollowerID int64 `json:"followerId"`
	FolloweeID int64 `json:"followeeId"`
}

type FollowUserReq struct {
	FollowerID int64 `json:"followerId"`
	FolloweeID int64 `json:"followeeId"`
}

type CancelFollowUserReq struct {
	FollowerID int64 `json:"followerId"`
	FolloweeID int64 `json:"followeeId"`
}

type GetFolloweeCountReq struct {
	UserID int64 `json:"userId"`
}

type GetFollowerCountReq struct {
	UserID int64 `json:"userId"`
}
