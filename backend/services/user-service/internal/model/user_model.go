package model

type User struct {
	ID    int32
	Name  string
	Email string
}

type Followers struct {
	ID int32
	Follower int32
	Followee int32
}

type Unfollow struct {
	ID int32
	Unfollower int32
	Unfollowee int32
}

// tlvz inutil isso ai