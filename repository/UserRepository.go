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
	//Get(uow *UnitOfWork, name int) (user *models.User, err error)
	GetLoggedInUser(uow *UnitOfWork, userId int) (*models.User, error)
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

//func (u userRepository) Get(uow *UnitOfWork, userId int) (user *models.User, err error) {
//
//	query := fmt.Sprintf(`
//									SELECT
//										id,
//										name,
//										password,
//										git_username
//									from
//										user
//									where
//										id =?`,
//	)
//	err = uow.Db.QueryRow(query, userId).Scan(&user.Id, &user.Name, &user.Password, &user.GitUsername)
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}

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

// deletes the token for the user and particular device
//func DeleteAuth(dbSvc *repository.DbSvc, userId int, authId string) error {
//
//	query, err := dbSvc.Db.Db.Prepare(`Delete * from authentication where user_id = ? and id = ? )`)
//	if err != nil {
//		// Todo error handling
//	}
//
//	_, err = query.Exec(userId, authId)
//	if err != nil {
//		// Todo error handling
//	}
//
//	return err
//}
