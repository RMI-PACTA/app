// Package db provides generalized utilities for interacting with some
// database, whether its an in-memory mock, a live SQL database, or something
// else.
package db

import (
	"errors"
	"fmt"

	"github.com/RMI/pacta"
)

type errNotFound struct {
	// id is the ID that wasn't found.
	id string
	// entityType is the type of the entity that the caller was looking for.
	entityType string
}

func (e *errNotFound) Error() string {
	return fmt.Sprintf("entity of type %q with ID %q was not found", e.entityType, e.id)
}

func NotFound[T ~string](id T, entityType string) error {
	return &errNotFound{id: string(id), entityType: entityType}
}

func (e *errNotFound) Is(target error) bool {
	_, ok := target.(*errNotFound)
	return ok
}

func IsNotFound(err error) bool {
	return errors.Is(err, &errNotFound{})
}

type Tx interface {
	Commit() error
	Rollback() error
}

type UpdateUserFn func(*pacta.User) error

func SetUserName(value string) UpdateUserFn {
	return func(u *pacta.User) error {
		u.Name = value
		return nil
	}
}

func SetUserEmail(value string) UpdateUserFn {
	return func(u *pacta.User) error {
		u.Email = value
		return nil
	}
}
