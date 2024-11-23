package services

import (
	"gin001/apis/models"
	"gin001/persistence/entity"
	"gin001/persistence/repository"
)

// UserService interface
type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService with repository
func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{userRepository: repository}
}

// GetUsers return array of UserEntity
func (s *UserService) GetUsers() ([]entity.UserEntity, error) {
	// implement the logic to get users from database or any other data source
	return []entity.UserEntity{}, nil
}

// GetUser with id return UserEntity
func (s *UserService) GetUser(id uint) (*models.UserInfo, error) {
	// implement the logic to get user by id from database or any other data source
	entity, err := s.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}
	return &models.UserInfo{
		ID:       entity.ID,
		Name:     entity.Name,
		BirthDay: entity.BirthDay,
		PhotoURL: entity.PhotoURL.String,
		Gender:   entity.Gender,
	}, nil
}

// CreateUser with UserEntity and return UserEntity
func (s *UserService) CreateUser(user entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to create a new user in database or any other data source
	return &user, nil
}

// UpdateUser with id, UserEntity, return UserEntity
func (s *UserService) UpdateUser(id uint, user entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to update an existing user in database or any other data source
	return &user, nil
}

// DeleteUser with id, if error will return
func (s *UserService) DeleteUser(id uint) error {
	// implement the logic to delete a user from database or any other data source
	return nil
}
