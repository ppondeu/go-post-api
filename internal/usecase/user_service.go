package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ppondeu/go-post-api/internal/dto"
	errs "github.com/ppondeu/go-post-api/internal/error"
	"github.com/ppondeu/go-post-api/internal/domain"
	"github.com/ppondeu/go-post-api/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserByID(ID uuid.UUID) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	CreateUser(createUserDto *dto.CreateUserDto) (*domain.User, error)
	UpdateUser(ID uuid.UUID, updateUserDto *dto.UpdateUserDto) (*domain.User, error)
	UpdateUserSession(userID uuid.UUID, refreshToken *string) error
	CreateUserAndSession(createUserDto *dto.CreateUserDto, refreshToken *string) (*domain.User, error)
	GetUserSession(userID uuid.UUID) (*domain.UserSession, error)
	DeleteUser(ID uuid.UUID) error
}

type UserServiceImpl struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &UserServiceImpl{userRepo}
}

func (s *UserServiceImpl) GetUserByID(ID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(createUserDto *dto.CreateUserDto) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserDto.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
		return nil, errs.NewInternalServerError()
	}
	user := &domain.User{
		Username: createUserDto.Username,
		Email:    createUserDto.Email,
		Password: string(hashedPassword),
		ShortBio: createUserDto.ShortBio,
	}
	result, err := s.userRepo.Create(user)
	if err != nil {
		logger.Error(err)
		return nil, errs.NewBadRequestError("Duplicate username or email")
	}
	return result, nil
}

func (s *UserServiceImpl) GetAllUsers() ([]domain.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (s *UserServiceImpl) GetUserByUsername(username string) (*domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		logger.Error(err)
		return nil, errs.NewBadRequestError("User with username not found")
	}
	return user, nil
}

func (s *UserServiceImpl) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		logger.Error(err)
		return nil, errs.NewBadRequestError("User with email not found")
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(ID uuid.UUID, updateUserDto *dto.UpdateUserDto) (*domain.User, error) {
	fmt.Println("updateUserDto: ", updateUserDto)
	user := &domain.User{
		Username: updateUserDto.Username,
		ShortBio: updateUserDto.ShortBio,
	}
	if updateUserDto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserDto.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error(err)
			return nil, errs.NewInternalServerError()
		}
		user.Password = string(hashedPassword)
	}

	result, err := s.userRepo.Update(ID, user)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *UserServiceImpl) DeleteUser(ID uuid.UUID) error {
	err := s.userRepo.Delete(ID)
	if err != nil {
		logger.Error(err)
		return errs.NewNotFoundError("User not found")
	}
	return nil
}

func (s *UserServiceImpl) CreateUserAndSession(createUserDto *dto.CreateUserDto, refreshToken *string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserDto.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
		return nil, errs.NewInternalServerError()
	}
	user := &domain.User{
		Username: createUserDto.Username,
		Email:    createUserDto.Email,
		Password: string(hashedPassword),
		ShortBio: createUserDto.ShortBio,
	}
	result, err := s.userRepo.CreateUserAndSession(user, refreshToken)
	if err != nil {
		logger.Error(err)
		return nil, errs.NewBadRequestError("Duplicate username or email")
	}
	return result, nil
}

func (s *UserServiceImpl) UpdateUserSession(userID uuid.UUID, refreshToken *string) error {
	err := s.userRepo.UpdateSession(userID, refreshToken)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (s *UserServiceImpl) GetUserSession(userID uuid.UUID) (*domain.UserSession, error) {
	session, err := s.userRepo.FindSession(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return session, nil
}
