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
	// implement the logic to create a new user in the database or storage system.
	insertSQL := "insert into users(name, birthday, gender, photo_url, email) values($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := u.dbCon.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}
	var id uint
	err = stmt.Get(&id, user.Name, user.BirthDay, user.Gender, user.PhotoURL, user.Email)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// CreateUser2 with userEntity
func (u *UserRepository) CreateUser2(user *entity.UserEntity) (*entity.UserEntity, error) {
	// implement the logic to create a new user in the database or storage system.
	insertSQL := "insert into users(name, birthday, gender, photo_url, email) values(:name, :birthday, :gender, :photo_url, :email) RETURNING id"
	stmt, err := u.dbCon.PrepareNamed(insertSQL)
	if err != nil {
		return nil, err
	}
	var id uint
	err = stmt.Get(&id, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// GetUser with ID
func (u *UserRepository) GetUser(id uint) (*entity.UserEntity, error) {
	// to retrieve a user from the database or storage system based on its ID.
	querySQL := "select * from users where id = $1"
	stmt, err := u.dbCon.Preparex(querySQL)
	if err != nil {
		return nil, err
	}
	var entity entity.UserEntity
	err = stmt.Get(&entity, id)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// UpdateUser with userEntity
func (u *UserRepository) UpdateUser(user *entity.UserEntity) (bool, error) {
	insertSQL := "update users set name = :name, birthday = :birthday, gender = :gender, photo_url = :photo_url, email = :email where id = :id"
	stmt, err := u.dbCon.PrepareNamed(insertSQL)
	if err != nil {
		return false, err
	}
	result := stmt.MustExec(user)
	count, err := result.RowsAffected()
	if count == 0 || err != nil {
		return false, err
	}
	return true, nil
}

// FindUsers with userEntity
func (u *UserRepository) FindUsers() ([]*entity.UserEntity, error) {
	//
	querySQL := "select id, name, birthday, gender, active, photo_url, email from users"
	stmt, err := u.dbCon.Preparex(querySQL)
	if err != nil {
		return nil, err
	}
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
