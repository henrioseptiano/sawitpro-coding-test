// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	RegisterUser(User) error
	CheckUser(context.Context, string) (*User, error)
	UpdateLoginUser(context.Context, User) error
	GetUserByUserId(context.Context, string) (*User, error)
	UpdateUserProfile(context.Context, User) error
	CheckPhoneNumber(context.Context, string) (int64, error)
}
