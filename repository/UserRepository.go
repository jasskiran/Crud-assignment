package repository

import (
	"fmt"

	"Assignment/models"
	_ "github.com/go-sql-driver/mysql"
)

type UserRepository interface {
	Add(uow *UnitOfWork, out *models.User) error
	Login(uow *UnitOfWork, name string) (*models.User, error)
	Update(uow *UnitOfWork, out *models.User, Id int) error
	GetLoggedInUser(uow *UnitOfWork, userId int) (*models.User, error)
	AddToken(uow *UnitOfWork, out *models.Auth) error
	GetToken(uow *UnitOfWork, userId int) (*models.Auth, error)
	DeleteToken(uow *UnitOfWork, userId int) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (u userRepository) Add(uow *UnitOfWork, out *models.User) error {
	query := fmt.Sprintf(`
		INSERT INTO user
		(
			id,
			name,
			password,
			git_username
		)
		VALUES
		(
			?,
			?,
			?,
			?
		)`,
	)

	stmt, err := uow.Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(out.Id, out.Name, out.Password, out.GitUsername)
	if err != nil {
		return err
	}
	return nil
}

func (u userRepository) Update(uow *UnitOfWork, user *models.User, Id int) error {

	query := fmt.Sprintf(`UPDATE 
										user 
									SET 
										name = ? 
									WHERE 
										id = ?`,
	)

	stmt, err := uow.Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Name, Id)
	if err != nil {
		return err
	}
	return nil
}

func (u userRepository) GetLoggedInUser(uow *UnitOfWork, userId int) (*models.User, error) {

	var user models.User
	var err error

	query := fmt.Sprintf(`
									SELECT
										id,
										name,
										password,
										git_username
											
									FROM
										user
									WHERE id = ?`,
	)
	err = uow.Db.QueryRow(query, userId).Scan(&user.Id, &user.Name, &user.Password, &user.GitUsername)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Login(uow *UnitOfWork, name string) (*models.User, error) {

	user := models.User{}
	var err error
	query := fmt.Sprintf(`
									SELECT
											id,
											name,
											password,
											git_username
									from
											user
									where
											name = ?`,
	)
	err = uow.Db.QueryRow(query, name).Scan(&user.Id, &user.Name, &user.Password, &user.GitUsername)
	return &user, err
}

func (u userRepository) AddToken(uow *UnitOfWork, out *models.Auth) error {
	query := fmt.Sprintf(`
		INSERT INTO authentication
		(
			id,
			user_id,
			token,
			active
		)
		VALUES
		(
			?,
			?,
			?,
			?
		)`,
	)

	stmt, err := uow.Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(out.Id, out.UserId, out.Token, true)
	if err != nil {
		return err
	}
	return nil
}

func (u userRepository) GetToken(uow *UnitOfWork, userId int) (*models.Auth, error) {

	var authentication models.Auth
	var err error
	query := fmt.Sprintf(`
									SELECT
										id,
										user_id,
										token,
										active
									from
										authentication
									where
										active = 1
									and user_id = ?`,
	)
	err = uow.Db.QueryRow(query, userId).Scan(&authentication.Id, &authentication.UserId, &authentication.Token, &authentication.Active)
	if err != nil {
		return nil, err
	}
	return &authentication, nil
}

func (u userRepository) DeleteToken(uow *UnitOfWork, userId int) error {
	query := fmt.Sprintf(`
		Update authentication
		set
			active = ?
		where 
			user_id = ?`,
	)

	stmt, err := uow.Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(false, userId)
	if err != nil {
		return err
	}
	return nil
}
