// Package db provides generalized utilities for interacting with some
// database, whether its an in-memory mock, a live SQL database, or something
// else.
package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/RMI/pacta/pacta"
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

func SetUserAdmin(value bool) UpdateUserFn {
	return func(u *pacta.User) error {
		u.Admin = value
		return nil
	}
}

func SetUserSuperAdmin(value bool) UpdateUserFn {
	return func(u *pacta.User) error {
		u.SuperAdmin = value
		return nil
	}
}

func SetUserName(value string) UpdateUserFn {
	return func(u *pacta.User) error {
		u.Name = value
		return nil
	}
}

func SetUserPreferredLanguage(value pacta.Language) UpdateUserFn {
	return func(u *pacta.User) error {
		u.PreferredLanguage = value
		return nil
	}
}

type UpdatePACTAVersionFn func(*pacta.PACTAVersion) error

func SetPACTAVersionName(value string) UpdatePACTAVersionFn {
	return func(v *pacta.PACTAVersion) error {
		v.Name = value
		return nil
	}
}

func SetPACTAVersionDescription(value string) UpdatePACTAVersionFn {
	return func(v *pacta.PACTAVersion) error {
		v.Description = value
		return nil
	}
}

func SetPACTAVersionDigest(value string) UpdatePACTAVersionFn {
	return func(v *pacta.PACTAVersion) error {
		v.Digest = value
		return nil
	}
}

type UpdateInitiativeFn func(*pacta.Initiative) error

func SetInitiativeName(value string) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.Name = value
		return nil
	}
}

func SetInitiativeAffiliation(value string) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.Affiliation = value
		return nil
	}
}

func SetInitiativePublicDescription(value string) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.PublicDescription = value
		return nil
	}
}

func SetInitiativeInternalDescription(value string) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.InternalDescription = value
		return nil
	}
}

func SetInitiativeRequiresInvitationToJoin(value bool) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.RequiresInvitationToJoin = value
		return nil
	}
}

func SetInitiativeIsAcceptingNewMembers(value bool) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.IsAcceptingNewMembers = value
		return nil
	}
}

func SetInitiativeIsAcceptingNewPortfolios(value bool) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.IsAcceptingNewPortfolios = value
		return nil
	}
}

func SetInitiativePACTAVersion(pvid pacta.PACTAVersionID) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.PACTAVersion = &pacta.PACTAVersion{ID: pvid}
		return nil
	}
}

func SetInitiativeLanguage(l pacta.Language) UpdateInitiativeFn {
	return func(v *pacta.Initiative) error {
		v.Language = l
		return nil
	}
}

type UpdateInitiativeInvitationFn func(*pacta.InitiativeInvitation) error

func SetInitiativeInvitationUsedAt(t time.Time) UpdateInitiativeInvitationFn {
	return func(ii *pacta.InitiativeInvitation) error {
		ii.UsedAt = t
		return nil
	}
}

func ClearInitiativeInvitationUsedBy() UpdateInitiativeInvitationFn {
	return func(ii *pacta.InitiativeInvitation) error {
		ii.UsedBy = nil
		return nil
	}
}

func SetInitiativeInvitationUsedBy(u pacta.UserID) UpdateInitiativeInvitationFn {
	return func(ii *pacta.InitiativeInvitation) error {
		if u == "" {
			return fmt.Errorf("cannot set used by to empty user ID")
		}
		ii.UsedBy = &pacta.User{ID: u}
		return nil
	}
}

type UpdateBlobFn func(*pacta.Blob) error

func SetBlobFileName(v string) UpdateBlobFn {
	return func(b *pacta.Blob) error {
		b.FileName = v
		return nil
	}
}

type UpdatePortfolioFn func(*pacta.Portfolio) error

func SetPortfolioName(value string) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Name = value
		return nil
	}
}

func SetPortfolioDescription(value string) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Description = value
		return nil
	}
}

func SetPortfolioHoldingsDate(value time.Time) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.HoldingsDate = value
		return nil
	}
}

func SetPortfolioOwner(value pacta.OwnerID) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Owner = &pacta.Owner{ID: value}
		return nil
	}
}

func SetPortfolioAdminDebugEnabled(value bool) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.AdminDebugEnabled = value
		return nil
	}
}

func SetPortfolioNumberOfRows(value int) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.NumberOfRows = value
		return nil
	}
}

type UpdatePortfolioGroupFn func(*pacta.PortfolioGroup) error

func SetPortfolioGroupName(value string) UpdatePortfolioGroupFn {
	return func(v *pacta.PortfolioGroup) error {
		v.Name = value
		return nil
	}
}

func SetPortfolioGroupDescription(value string) UpdatePortfolioGroupFn {
	return func(v *pacta.PortfolioGroup) error {
		v.Description = value
		return nil
	}
}

func SetPortolioGroupOwner(value pacta.OwnerID) UpdatePortfolioGroupFn {
	return func(v *pacta.PortfolioGroup) error {
		v.Owner = &pacta.Owner{ID: value}
		return nil
	}
}
