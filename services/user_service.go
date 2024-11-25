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

func (s *UserService) mapToModel(ue *entity.UserEntity) *models.UserInfo {
	return &models.UserInfo{
		ID:        ue.ID,
		Name:      ue.Name,
		BirthDay:  ue.BirthDay,
		PhotoURL:  ue.PhotoURL,
		Gender:    ue.Gender,
		Email:     ue.Email,
		CreatedAt: ue.CreatedAt,
		UpdatedAt: ue.UpdatedAt,
	}
}

// GetUsers return array of UserInfo
func (s *UserService) GetUsers() ([]*models.UserInfo, error) {
	// implement the logic to get users from database or any other data source
	list, err := s.userRepository.FindUsers()
	if err != nil {
		return nil, err
	}
	var result []*models.UserInfo
	for _, ue := range list {
		result = append(result, s.mapToModel(ue))
	}
	return result, nil
}

// GetUser with id return UserInfo
func (s *UserService) GetUser(id uint) (*models.UserInfo, error) {
	// implement the logic to get user by id from database or any other data source
	ue, err := s.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}
	return s.mapToModel(ue), nil
}

// CreateUser with UserEntity and return UserInfo
func (s *UserService) CreateUser(signUp *models.UserSignUp) (*models.UserInfo, error) {
	// implement the logic to create a new user in database or any other data source
	ue := entity.UserEntity{
		Name:     signUp.Name,
		PhotoURL: signUp.PhotoURL,
		BirthDay: signUp.BirthDay,
		Gender:   signUp.Gender,
		Email:    signUp.Email,
	}
	// ueCreated, err := s.userRepository.CreateUser(&ue)
	ueCreated, err := s.userRepository.CreateUser2(&ue)
	if err != nil {
		return nil, err
	}
	return s.mapToModel(ueCreated), nil
}

// UpdateUser with id, UserInfoUpdate, return UserInfo
func (s *UserService) UpdateUser(id uint, userInfoUpdate *models.UserInfoUpdate) (*models.UserInfo, error) {
	// implement the logic to update an existing user in database or any other data source
	ue := entity.UserEntity{
		ID:       id,
		Name:     userInfoUpdate.Name,
		PhotoURL: userInfoUpdate.PhotoURL,
		BirthDay: userInfoUpdate.BirthDay,
		Gender:   userInfoUpdate.Gender,
		Email:    userInfoUpdate.Email,
	}
	updated, err := s.userRepository.UpdateUser(&ue)
	if !updated || err != nil {
		return nil, err
	}
	ueUpdated, err := s.GetUser(id)
	if err != nil {
		return nil, err
	}
	return ueUpdated, nil
}

// DeleteUser with id, if error will return
func (s *UserService) DeleteUser(id uint) error {
	// implement the logic to delete a user from database or any other data source
	return nil
}
