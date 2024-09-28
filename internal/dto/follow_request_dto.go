package dto

type FollowRequestDto struct {
	FollowedID string `json:"followedID" validate:"required,uuid"`
	FollowerID string `json:"followerID" validate:"required,uuid"`
}
