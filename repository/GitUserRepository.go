package repository

import (
	"fmt"

	"Assignment/models"
)

type GitUserRepository interface {
	CreateGitUser(uow *UnitOfWork, out *models.Github) error
	GetGitUser(uow *UnitOfWork, userId int) (gitUser *models.Github, err error)
}


type gitUserRepository struct {

}

func NewGitUserRepository() GitUserRepository{
	return &gitUserRepository{}
}

func (u *gitUserRepository) CreateGitUser(uow *UnitOfWork, out *models.Github)  error{

	query := fmt.Sprintf(`
		INSERT INTO github
		(
			id,
			username
		)
		VALUES
		(
			?,
			?
		)`,
	)

	stmt , err := uow.Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(out.Id, out.UserName)
	if err != nil {
		return err

	}
	return nil
}

func (u *gitUserRepository) GetGitUser(uow *UnitOfWork, userId int) (gitUser *models.Github, err error) {

	query := fmt.Sprintf(`
									SELECT
										github.id,
										github.username
									FROM
										github
									LEFT OUTER JOIN user ON github.username = user.git_username
									WHERE user.id = ?`,
	)
	err = uow.Db.QueryRow(query, userId).Scan(&gitUser.Id, &gitUser.UserName)
	if err != nil {
		return
	}
	return
}
