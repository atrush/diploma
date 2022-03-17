package storage

import (
	"context"
	"github.com/atrush/diploma.git/model"
)

//  Storage is the interface that wraps methods for working with the database.
type Storage interface {

	//  User returns repository for working with users.
	User() UserRepository
	//  Order returns repository for working with orders.
	Order() OrderRepository
	//  Close closes storage connection.
	Close()
}

type UserRepository interface {
	//  Adds new user to storage
	Create(ctx context.Context, login string, hash string) (model.User, error)
	//  Returns user from storage
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type OrderRepository interface {
}
