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

func SetPortfolioPropertyHoldingsDate(value *pacta.HoldingsDate) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Properties.HoldingsDate = value
		return nil
	}
}

func SetPortfolioPropertyESG(value *bool) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Properties.ESG = value
		return nil
	}
}

func SetPortfolioPropertyExternal(value *bool) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Properties.External = value
		return nil
	}
}

func SetPortfolioPropertyEngagementStrategy(value *bool) UpdatePortfolioFn {
	return func(v *pacta.Portfolio) error {
		v.Properties.EngagementStrategy = value
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

func SetPortfolioGroupOwner(value pacta.OwnerID) UpdatePortfolioGroupFn {
	return func(v *pacta.PortfolioGroup) error {
		v.Owner = &pacta.Owner{ID: value}
		return nil
	}
}

type UpdateAnalysisFn func(*pacta.Analysis) error

func SetAnalysisOwner(value pacta.OwnerID) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.Owner = &pacta.Owner{ID: value}
		return nil
	}
}

func SetAnalysisName(value string) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.Name = value
		return nil
	}
}

func SetAnalysisDescription(value string) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.Description = value
		return nil
	}
}

func SetAnalysisRanAt(value time.Time) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.RanAt = value
		return nil
	}
}

func SetAnalysisCompletedAt(value time.Time) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.CompletedAt = value
		return nil
	}
}

func SetAnalysisFailureCode(value pacta.FailureCode) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.FailureCode = value
		return nil
	}
}

func SetAnalysisFailureMessage(value string) UpdateAnalysisFn {
	return func(v *pacta.Analysis) error {
		v.FailureMessage = value
		return nil
	}
}

type UpdateAnalysisArtifactFn func(*pacta.AnalysisArtifact) error

func SetAnalysisArtifactAdminDebugEnabled(value bool) UpdateAnalysisArtifactFn {
	return func(v *pacta.AnalysisArtifact) error {
		v.AdminDebugEnabled = value
		return nil
	}
}

func SetAnalysisArtifactSharedToPublic(value bool) UpdateAnalysisArtifactFn {
	return func(v *pacta.AnalysisArtifact) error {
		v.SharedToPublic = value
		return nil
	}
}

type UpdateIncompleteUploadFn func(*pacta.IncompleteUpload) error

func SetIncompleteUploadOwner(value pacta.OwnerID) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Owner = &pacta.Owner{ID: value}
		return nil
	}
}

func SetIncompleteUploadAdminDebugEnabled(value bool) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.AdminDebugEnabled = value
		return nil
	}
}

func SetIncompleteUploadName(value string) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Name = value
		return nil
	}
}

func SetIncompleteUploadDescription(value string) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Description = value
		return nil
	}
}

func SetIncompleteUploadRanAt(value time.Time) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.RanAt = value
		return nil
	}
}

func SetIncompleteUploadCompletedAt(value time.Time) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.CompletedAt = value
		return nil
	}
}

func SetIncompleteUploadFailureCode(value pacta.FailureCode) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.FailureCode = value
		return nil
	}
}

func SetIncompleteUploadFailureMessage(value string) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.FailureMessage = value
		return nil
	}
}

func SetIncompleteUploadPropertyHoldingsDate(value *pacta.HoldingsDate) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Properties.HoldingsDate = value
		return nil
	}
}

func SetIncompleteUploadPropertyESG(value *bool) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Properties.ESG = value
		return nil
	}
}

func SetIncompleteUploadPropertyExternal(value *bool) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Properties.External = value
		return nil
	}
}

func SetIncompleteUploadPropertyEngagementStrategy(value *bool) UpdateIncompleteUploadFn {
	return func(v *pacta.IncompleteUpload) error {
		v.Properties.EngagementStrategy = value
		return nil
	}
}

type UpdateInitiativeUserRelationshipFn func(*pacta.InitiativeUserRelationship) error

func SetInitiativeUserRelationshipMember(value bool) UpdateInitiativeUserRelationshipFn {
	return func(v *pacta.InitiativeUserRelationship) error {
		v.Member = value
		return nil
	}
}

func SetInitiativeUserRelationshipManager(value bool) UpdateInitiativeUserRelationshipFn {
	return func(v *pacta.InitiativeUserRelationship) error {
		v.Manager = value
		return nil
	}
}
