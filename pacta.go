// Package pacta contains domain types for the PACTA ecosystem
package pacta

import "time"

type Provider string

const (
	EmailAndPass Provider = "EMAIL_AND_PASS"
	Facebook     Provider = "FACEBOOK"
	Google       Provider = "GOOGLE"
)

type UserID string

// This is just an example struct
type User struct {
	ID                UserID
	Name              string
	Email             string
	CreatedAt         time.Time
	AuthnProviderType Provider
	AuthnProviderID   UserID
}
