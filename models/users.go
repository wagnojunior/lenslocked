package models

import "database/sql"

type User struct {
	ID           int
	Email        string
	PasswordHash int
}

// UserService is a type that has the database connection as a field.
// There are two main reasons why this approch is preferred:
// 1) decouple the methods associated with this type (i.e.: Create method)
// from the underlying database implementation. In order words, no prior
// knowledge of how the database was implemented is necessary for the
// associated methods to be called.
// 2) possibility to use interface (don't fully understand this one yet)
type UserService struct {
	DB *sql.DB
}
