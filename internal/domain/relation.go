package domain

type Relation struct {
	FolloweeId int64
	FollowerId int64
}

type RelationStats struct {
	FollowerCount int64
	FolloweeCount int64
}
