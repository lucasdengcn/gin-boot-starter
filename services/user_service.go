package services

import (
	"gin-boot-starter/apis/models"
	"gin-boot-starter/core/exception"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/core/security"
	"gin-boot-starter/persistence/entity"
	"gin-boot-starter/persistence/repository"
	"sync"

	"github.com/gin-gonic/gin"
)

var onceUserService sync.Once
var instanceUserService *UserService

// UserService interface
type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService with repository
func NewUserService(repository *repository.UserRepository) *UserService {
	onceUserService.Do(func() {
		instanceUserService = &UserService{userRepository: repository}
	})
	return instanceUserService
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
func (s *UserService) GetUsers(ctx *gin.Context) []*models.UserInfo {
	// implement the logic to get users from database or any other data source
	list, err := s.userRepository.FindUsers(ctx)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("UserRepository.FindUsers error.")
		panic(exception.NewRepositoryError(ctx, "USER_GETS_500", err.Error()))
	}
	var result []*models.UserInfo
	for _, ue := range list {
		result = append(result, s.mapToModel(ue))
	}
	return result
}

// GetUser with id return UserInfo
func (s *UserService) GetUser(ctx *gin.Context, id uint) *models.UserInfo {
	// implement the logic to get user by id from database or any other data source
	logging.Debug(ctx).Msgf("GetUser with id:%v", id)
	ue, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("UserRepository.GetUser error. id=%v", id)
		panic(exception.NewRepositoryError(ctx, "USER_GET_500", err.Error()))
	}
	if nil == ue {
		panic(exception.NewEntityNotFoundError(ctx, id, "User not found"))
	}
	return s.mapToModel(ue)
}

// CreateUser with UserEntity and return UserInfo
func (s *UserService) CreateUser(ctx *gin.Context, signUp *models.UserSignUp) *models.UserInfo {
	// implement the logic to create a new user in database or any other data source
	hashPassword, err := security.HashPassword(signUp.Password)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("UserService.CreateUser error. data: %v", signUp)
		panic(exception.NewServiceError(ctx, "USER_PASSWORD_HASH_500", err.Error()))
	}
	ue := entity.UserEntity{
		Name:     signUp.Name,
		PhotoURL: signUp.PhotoURL,
		BirthDay: signUp.BirthDay,
		Gender:   signUp.Gender,
		Email:    signUp.Email,
		Roles:    security.RoleUser,
		Password: hashPassword,
	}
	ueCreated, err := s.userRepository.CreateUser2(ctx, &ue)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("UserRepository.CreateUser error. data: %v", signUp)
		panic(exception.NewRepositoryError(ctx, "USER_CREATE_500", err.Error()))
	}
	return s.mapToModel(ueCreated)
}

// UpdateUser with id, UserInfoUpdate, return UserInfo
func (s *UserService) UpdateUser(ctx *gin.Context, id uint, userInfoUpdate *models.UserInfoUpdate) *models.UserInfo {
	// implement the logic to update an existing user in database or any other data source
	ue := entity.UserEntity{
		ID:       id,
		Name:     userInfoUpdate.Name,
		PhotoURL: userInfoUpdate.PhotoURL,
		BirthDay: userInfoUpdate.BirthDay,
		Gender:   userInfoUpdate.Gender,
		Email:    userInfoUpdate.Email,
	}
	updated, err := s.userRepository.UpdateUser(ctx, &ue)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("UserRepository.UpdateUser error. data: %v", userInfoUpdate)
		panic(exception.NewRepositoryError(ctx, "USER_UPDATE_500", err.Error()))
	}
	if !updated {
		panic(exception.NewEntityNotFoundError(ctx, id, "User not found."))
	}
	// for testing transaction rollback, if panic here, UpdateUser will rollback.
	// panic(core.NewServiceError(400, "Something is wrong"))
	ueUpdated := s.GetUser(ctx, id)
	return ueUpdated
}

// DeleteUser with id, if error will return
func (s *UserService) DeleteUser(ctx *gin.Context, id uint) error {
	// implement the logic to delete a user from database or any other data source
	return nil
}

// VerifyPassword sign in user and verify password
func (s *UserService) VerifyPassword(ctx *gin.Context, signIn *models.UserSignIn) *models.UserInfo {
	// implement the logic to sign in user in database or any other data source
	ue, err := s.userRepository.GetUserByEmail(ctx, signIn.Email)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("UserRepository.GetUserByEmail error. email=%v", signIn.Email)
		panic(exception.NewRepositoryError(ctx, "USER_GET_BY_EMAIL_500", err.Error()))
	}
	if ue == nil {
		logging.Error(ctx).Err(err).Msgf("UserService.VerifyPassword invalid email. email=%v", signIn.Email)
		panic(exception.NewAuthError(ctx, "USER_ACCOUNT_404", "Invalid password or email"))
	}
	ok := security.VerifyPassword(signIn.Password, ue.Password)
	if !ok {
		logging.Error(ctx).Err(err).Msgf("UserService.VerifyPassword invalid password. email=%v", signIn.Email)
		panic(exception.NewAuthError(ctx, "USER_ACCOUNT_400", "Invalid password or email"))
	}
	return s.mapToModel(ue)
}
