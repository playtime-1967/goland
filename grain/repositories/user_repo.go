package repositories

import (
	"grain/entities"

	"github.com/gocql/gocql"
)

type UserRepo interface {
	Create(user entities.User) error
}

type CassandraUserRepo struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) UserRepo {
	return &CassandraUserRepo{session: session}
}

func (r *CassandraUserRepo) Create(user entities.User) error {
	return r.session.Query(`INSERT INTO users (id, name, email) VALUES (?, ?, ?)`, user.ID, user.Name, user.Email).Exec()
}
