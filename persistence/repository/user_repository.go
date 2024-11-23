package repository

import (
	"gin001/persistence/entity"

	"github.com/jmoiron/sqlx"
)

// UserRepository interface
type UserRepository struct {
	dbCon *sqlx.DB
}

// NewUserRepository with DB connection
func NewUserRepository(dbCon *sqlx.DB) *UserRepository {
	return &UserRepository{dbCon: dbCon}
}

// CreateUser with userEntity
func (u *UserRepository) CreateUser(user *entity.UserEntity) (*entity.UserEntity, error) {
	// TO DO: implement the logic to create a new user in the database or storage system.
	return &entity.UserEntity{}, nil
}

// GetUser with ID
func (u *UserRepository) GetUser(id uint) (*entity.UserEntity, error) {
	// TO DO: implement the logic to retrieve a user from the database or storage system based on its ID.
	row := u.dbCon.QueryRowx("select id, name, birthday, gender, active, photo_url from users where id = $1", id)
	var entity entity.UserEntity
	err := row.StructScan(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// UpdateUser with userEntity
func (u *UserRepository) UpdateUser(user *entity.UserEntity) (*entity.UserEntity, error) {
	// TO DO: implement the logic to update an existing user in the database or storage system.
	return &entity.UserEntity{}, nil
}
