package repository

import "database/sql"

type UnitOfWork struct {
	Db     *sql.DB
}

func NewUnitOfWork(db  *sql.DB) *UnitOfWork{
	return &UnitOfWork{Db: db}
}