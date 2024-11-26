package services

import (
	"gin001/apis/models"
	"gin001/core"
	"gin001/core/logging"
	"gin001/persistence/entity"
	"gin001/persistence/repository"

	"github.com/gin-gonic/gin"
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
func (s *UserService) GetUsers(c *gin.Context) []*models.UserInfo {
	// implement the logic to get users from database or any other data source
	list, err := s.userRepository.FindUsers(c)
	if err != nil {
		logging.Error(c).Err(err).Msgf("UserRepository.FindUsers error.")
		panic(core.NewServiceError(500, "FindUsers error", "UserRepository"))
	}
	var result []*models.UserInfo
	for _, ue := range list {
		result = append(result, s.mapToModel(ue))
	}
	return result
}

// GetUser with id return UserInfo
func (s *UserService) GetUser(c *gin.Context, id uint) *models.UserInfo {
	// implement the logic to get user by id from database or any other data source
	logging.Debug(c).Msgf("GetUser with id:%v", id)
	ue, err := s.userRepository.GetUser(c, id)
	if err != nil {
		logging.Error(c).Err(err).Msgf("UserRepository.GetUser error. id=%v", id)
		panic(core.NewServiceError(500, "GetUser error", "UserRepository"))
	}
	if nil == ue {
		panic(core.NewEntityNotFoundError(id, "GetUser not found", "User"))
	}
	return s.mapToModel(ue)
}

// CreateUser with UserEntity and return UserInfo
func (s *UserService) CreateUser(c *gin.Context, signUp *models.UserSignUp) *models.UserInfo {
	// implement the logic to create a new user in database or any other data source
	ue := entity.UserEntity{
		Name:     signUp.Name,
		PhotoURL: signUp.PhotoURL,
		BirthDay: signUp.BirthDay,
		Gender:   signUp.Gender,
		Email:    signUp.Email,
	}
	// ueCreated, err := s.userRepository.CreateUser(&ue)
	ueCreated, err := s.userRepository.CreateUser2(c, &ue)
	if err != nil {
		logging.Error(c).Err(err).Msgf("UserRepository.CreateUser error. data: %v", signUp)
		panic(core.NewServiceError(500, "CreateUser error", "UserRepository"))
	}
	return s.mapToModel(ueCreated)
}

// UpdateUser with id, UserInfoUpdate, return UserInfo
func (s *UserService) UpdateUser(c *gin.Context, id uint, userInfoUpdate *models.UserInfoUpdate) *models.UserInfo {
	// implement the logic to update an existing user in database or any other data source
	ue := entity.UserEntity{
		ID:       id,
		Name:     userInfoUpdate.Name,
		PhotoURL: userInfoUpdate.PhotoURL,
		BirthDay: userInfoUpdate.BirthDay,
		Gender:   userInfoUpdate.Gender,
		Email:    userInfoUpdate.Email,
	}
	updated, err := s.userRepository.UpdateUser(c, &ue)
	if err != nil {
		panic(err)
	}
	if !updated {
		panic(core.NewEntityNotFoundError(id, "Update User not found.", "User"))
	}
	// for testing transaction rollback, if panic here, UpdateUser will rollback.
	// panic(core.NewServiceError(400, "Something is wrong"))
	ueUpdated := s.GetUser(c, id)
	return ueUpdated
}

// DeleteUser with id, if error will return
func (s *UserService) DeleteUser(c *gin.Context, id uint) error {
	// implement the logic to delete a user from database or any other data source
	return nil
}
