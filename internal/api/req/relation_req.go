package req

type ListFollowerRelationsReq struct {
	FollowerID int64  `json:"followerId"`     // 关注者
	Page       int    `json:"page,omitempty"` // 当前页码
	Size       *int64 `json:"size,omitempty"` // 每页数据量
}

type ListFolloweeRelationsReq struct {
	FolloweeID int64  `json:"followeeId"`     // 被关注者
	Page       int    `json:"page,omitempty"` // 当前页码
	Size       *int64 `json:"size,omitempty"` // 每页数据量
}

type GetRelationInfoReq struct {
	FollowerID int64 `json:"followerId"` // 关注者
	FolloweeID int64 `json:"followeeId"` // 被关注者
}

type FollowUserReq struct {
	FollowerID int64 `json:"followerId"` // 关注者
	FolloweeID int64 `json:"followeeId"` // 被关注者
}

type CancelFollowUserReq struct {
	FollowerID int64 `json:"followerId"` // 关注者
	FolloweeID int64 `json:"followeeId"` // 被关注者
}

type GetFolloweeCountReq struct {
	UserID int64 `json:"userId"`
}

type GetFollowerCountReq struct {
	UserID int64 `json:"userId"`
}
