package repository

import (
	"gin-boot-starter/core/logging"
	"gin-boot-starter/persistence/entity"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var onceUserRepository sync.Once
var instanceUserRepository *UserRepository

// UserRepository interface
type UserRepository struct {
	TransactionRepo
}

// NewUserRepository with DB connection
func NewUserRepository(dbCon *sqlx.DB) *UserRepository {
	onceUserRepository.Do(func() {
		instanceUserRepository = &UserRepository{
			TransactionRepo: NewTransactionRepo(dbCon),
		}
	})
	return instanceUserRepository
}

// CreateUser with userEntity
func (u *UserRepository) CreateUser(ctx *gin.Context, user *entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to create a new user in the database or storage system.
	insertSQL := "insert into users(name, birthday, gender, photo_url, email, hashed_password, roles) values($1, $2, $3, $4, $5, $6, %7) RETURNING id"
	stmt := u.prepareStatement(ctx, insertSQL)
	var id uint
	err := stmt.Get(&id, user.Name, user.BirthDay, user.Gender, user.PhotoURL, user.Email, user.Password, user.Roles)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// CreateUser2 with userEntity
func (u *UserRepository) CreateUser2(ctx *gin.Context, user *entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to create a new user in the database or storage system.
	insertSQL := "insert into users(name, birthday, gender, photo_url, email, hashed_password, roles) values(:name, :birthday, :gender, :photo_url, :email, :hashed_password, :roles) RETURNING id"
	stmt := u.prepareNamed(ctx, insertSQL)
	var id uint
	err := stmt.Get(&id, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// GetUser with ID
func (u *UserRepository) GetUser(ctx *gin.Context, id uint) (*entity.UserEntity, error) {
	// to retrieve a user from the database or storage system based on its ID.
	logging.Debug(ctx).Msgf("GetUser with id:%v", id)
	querySQL := "select * from users where id = $1"
	stmt := u.prepareStatement(ctx, querySQL)
	var entity entity.UserEntity
	err := stmt.Get(&entity, id)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetUserByEmail with Email
func (u *UserRepository) GetUserByEmail(ctx *gin.Context, email string) (*entity.UserEntity, error) {
	// to retrieve a user from the database or storage system based on its ID.
	logging.Debug(ctx).Msgf("GetUserByEmail with email:%v", email)
	querySQL := "select * from users where email = $1"
	stmt := u.prepareStatement(ctx, querySQL)
	var entity entity.UserEntity
	err := stmt.Get(&entity, email)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// UpdateUser with userEntity
func (u *UserRepository) UpdateUser(ctx *gin.Context, user *entity.UserEntity) (bool, error) {
	sql := "update users set name = :name, birthday = :birthday, gender = :gender, photo_url = :photo_url, email = :email where id = :id"
	stmt := u.prepareNamed(ctx, sql)
	result := stmt.MustExec(user)
	count, err := result.RowsAffected()
	if count == 0 || err != nil {
		return false, err
	}
	return true, nil
}

// FindUsers with userEntity
func (u *UserRepository) FindUsers(ctx *gin.Context) ([]*entity.UserEntity, error) {
	//
	querySQL := "select id, name, birthday, gender, active, photo_url, email from users"
	stmt := u.prepareStatement(ctx, querySQL)
	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}

	var list []*entity.UserEntity
	for rows.Next() {
		var e entity.UserEntity
		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	return list, nil
}
