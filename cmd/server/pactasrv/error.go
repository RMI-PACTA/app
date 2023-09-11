package pactasrv

import (
	"fmt"

	api "github.com/RMI/pacta/openapi/pacta"
)

type errBadRequest struct {
	Field   string
	Message string
}

func (e *errBadRequest) Code() int32 {
	return 400
}

func (e *errBadRequest) Error() string {
	return fmt.Sprintf("bad request: field %q: %s", e.Field, e.Message)
}

func (e *errBadRequest) Is(target error) bool {
	_, ok := target.(*errBadRequest)
	return ok
}

func errorBadRequest(field string, message string) error {
	return &errBadRequest{Field: field, Message: message}
}

type errUnauthorized struct {
	Action   string
	Resource string
}

func (e *errUnauthorized) Code() int32 {
	return 401
}

func (e *errUnauthorized) Error() string {
	return fmt.Sprintf("unauthorized to %s %s", e.Action, e.Resource)
}

func (e *errUnauthorized) Is(target error) bool {
	_, ok := target.(*errUnauthorized)
	return ok
}

func errorUnauthorized(action string, resource string) error {
	return &errUnauthorized{Action: action, Resource: resource}
}

type errForbidden struct {
	Action   string
	Resource string
}

func (e *errForbidden) Code() int32 {
	return 403
}

func (e *errForbidden) Error() string {
	return fmt.Sprintf("user is not allowed to %s %s", e.Action, e.Resource)
}

func (e *errForbidden) Is(target error) bool {
	_, ok := target.(*errForbidden)
	return ok
}

func errorForbidden(action string, resource string) error {
	return &errForbidden{Action: action, Resource: resource}
}

type errNotFound struct {
	What string
	With string
}

func (e *errNotFound) Code() int32 {
	return 404
}

func (e *errNotFound) Error() string {
	return fmt.Sprintf("not found: %s with %s", e.What, e.With)
}

func (e *errNotFound) Is(target error) bool {
	_, ok := target.(*errNotFound)
	return ok
}

func errorNotFound(what string, with string) error {
	return &errNotFound{What: what, With: with}
}

type errInternal struct {
	What string
}

func (e *errInternal) Code() int32 {
	return 500
}

func (e *errInternal) Error() string {
	return fmt.Sprintf("internal error: %s", e.What)
}

func (e *errInternal) Is(target error) bool {
	_, ok := target.(*errInternal)
	return ok
}

func errorInternal(err error) error {
	return &errInternal{What: err.Error()}
}

type errNotImplemented struct {
	What string
}

func (e *errNotImplemented) Code() int32 {
	return 501
}

func (e *errNotImplemented) Error() string {
	return fmt.Sprintf("not implemented: %s", e.What)
}

func (e *errNotImplemented) Is(target error) bool {
	_, ok := target.(*errNotImplemented)
	return ok
}

func errorNotImplemented(what string) error {
	return &errNotImplemented{What: what}
}

func errToAPIError(err error) api.Error {
	if e, ok := err.(*errBadRequest); ok {
		return api.Error{Code: e.Code(), Message: e.Error()}
	}
	if e, ok := err.(*errUnauthorized); ok {
		return api.Error{Code: e.Code(), Message: e.Error()}
	}
	if e, ok := err.(*errForbidden); ok {
		return api.Error{Code: e.Code(), Message: e.Error()}
	}
	if e, ok := err.(*errNotFound); ok {
		return api.Error{Code: e.Code(), Message: e.Error()}
	}
	if e, ok := err.(*errInternal); ok {
		return api.Error{Code: e.Code(), Message: e.Error()}
	}
	if e, ok := err.(*errNotImplemented); ok {
		return api.Error{Code: e.Code(), Message: e.Error()}
	}
	// TODO: log here for an unexpected error condition.
	return api.Error{
		Code:    500,
		Message: "an unexpected error occurred",
	}
}
