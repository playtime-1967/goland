package entities

import "github.com/gocql/gocql"

type User struct {
	ID    gocql.UUID `json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
}

func NewUser(name, email string) User {

	return User{
		Name:  name,
		Email: email,
	}
}
