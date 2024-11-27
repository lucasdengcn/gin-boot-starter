package repository

import (
	"context"
	"gin001/persistence/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// UserRepository interface
type UserRepository struct {
	TransactionRepo
}

// NewUserRepository with DB connection
func NewUserRepository(dbCon *sqlx.DB) *UserRepository {
	return &UserRepository{
		TransactionRepo: NewTransactionRepo(dbCon),
	}
}

// CreateUser with userEntity
func (u *UserRepository) CreateUser(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to create a new user in the database or storage system.
	insertSQL := "insert into users(name, birthday, gender, photo_url, email) values($1, $2, $3, $4, $5) RETURNING id"
	stmt := u.prepareStatement(ctx, insertSQL)
	var id uint
	err := stmt.Get(&id, user.Name, user.BirthDay, user.Gender, user.PhotoURL, user.Email)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// CreateUser2 with userEntity
func (u *UserRepository) CreateUser2(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to create a new user in the database or storage system.
	insertSQL := "insert into users(name, birthday, gender, photo_url, email) values(:name, :birthday, :gender, :photo_url, :email) RETURNING id"
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
func (u *UserRepository) GetUser(ctx context.Context, id uint) (*entity.UserEntity, error) {
	// to retrieve a user from the database or storage system based on its ID.
	log.Debug().Msgf("GetUser with id:%v", id)
	querySQL := "select * from users where id = $1"
	stmt := u.prepareStatement(ctx, querySQL)
	var entity entity.UserEntity
	err := stmt.Get(&entity, id)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// UpdateUser with userEntity
func (u *UserRepository) UpdateUser(ctx context.Context, user *entity.UserEntity) (bool, error) {
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
func (u *UserRepository) FindUsers(ctx context.Context) ([]*entity.UserEntity, error) {
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
