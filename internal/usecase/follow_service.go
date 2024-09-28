package usecase

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/domain"
	"github.com/ppondeu/go-post-api/internal/dto"
	"github.com/ppondeu/go-post-api/internal/errors"
	"github.com/ppondeu/go-post-api/internal/logger"
)

type FollowService interface {
	Follow(followerID, followedID uuid.UUID) error
	Unfollow(followerID, followedID uuid.UUID) error
	GetFollowers(userID uuid.UUID) ([]dto.UserResponseDto, error)
	GetFollowedUsers(userID uuid.UUID) ([]dto.UserResponseDto, error)
}

type followServiceImpl struct {
	followRepo  domain.FollowRepository
	userService UserService
}

func NewFollowService(followRepo domain.FollowRepository, userService UserService) FollowService {
	return &followServiceImpl{
		followRepo:  followRepo,
		userService: userService,
	}
}

func (s *followServiceImpl) Follow(followerID, followedID uuid.UUID) error {
	if followerID == followedID {
		logger.Error("cannot follow yourself")
		return errors.NewBadRequestError("cannot follow yourself")
	}

	follow, err := s.followRepo.FindByFollowerIDAndFollowedID(followerID, followedID)
	if err == nil && follow != nil {
		logger.Error("already followed")
		return nil
	}

	follow = &domain.Follow{
		FollowerID: followerID.String(),
		FollowedID: followedID.String(),
	}

	_, err = s.followRepo.Create(follow)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (s *followServiceImpl) Unfollow(followerID, followedID uuid.UUID) error {
	if followerID == followedID {
		logger.Error("cannot unfollow yourself")
		return errors.NewBadRequestError("cannot unfollow yourself")
	}

	follow, err := s.followRepo.FindByFollowerIDAndFollowedID(followerID, followedID)
	if err != nil {
		logger.Error(err)
		return errors.NewNotFoundError("not followed")
	}
	if follow == nil {
		logger.Error("not followed")
		return errors.NewNotFoundError("not followed")
	}

	ID, err := uuid.Parse(follow.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.followRepo.Delete(ID)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (s *followServiceImpl) GetFollowers(userID uuid.UUID) ([]dto.UserResponseDto, error) {
	follows, err := s.followRepo.FindFollowersByUserID(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var userResponseDtos []dto.UserResponseDto
	for _, follow := range follows {
		ID, err := uuid.Parse(follow.FollowerID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		user, err := s.userService.GetUserByID(ID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		userResponseDtos = append(userResponseDtos, dto.UserResponseDto{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}
	return userResponseDtos, nil
}

func (s *followServiceImpl) GetFollowedUsers(userID uuid.UUID) ([]dto.UserResponseDto, error) {
	follows, err := s.followRepo.FindFollowedUsersByUserID(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var userResponseDtos []dto.UserResponseDto
	for _, follow := range follows {
		ID, err := uuid.Parse(follow.FollowedID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		user, err := s.userService.GetUserByID(ID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		userResponseDtos = append(userResponseDtos, dto.UserResponseDto{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}
	return userResponseDtos, nil
}
