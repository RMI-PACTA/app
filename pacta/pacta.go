// Package pacta contains domain types for the PACTA ecosystem
package pacta

import (
	"fmt"
	"strings"
	"time"
)

type AuthnMechanism string

const (
	AuthnMechanism_EmailAndPass AuthnMechanism = "EMAIL_AND_PASS"
)

var AuthnMechanismValues = []AuthnMechanism{
	AuthnMechanism_EmailAndPass,
}

func ParseAuthnMechanism(s string) (AuthnMechanism, error) {
	if s == "EMAIL_AND_PASS" {
		return AuthnMechanism_EmailAndPass, nil
	}
	return "", fmt.Errorf("unknown AuthnMechanism: %q", s)
}

type Language string

const (
	Language_EN          Language = "en"
	Language_DE          Language = "de"
	Language_FR          Language = "fr"
	Language_ES          Language = "es"
	Language_Unspecified Language = "unspecified"
)

var LanguageValues = []Language{
	Language_EN,
	Language_DE,
	Language_FR,
	Language_ES,
}

func ParseLanguage(s string) (Language, error) {
	switch s {
	case "en":
		return Language_EN, nil
	case "de":
		return Language_DE, nil
	case "fr":
		return Language_FR, nil
	case "es":
		return Language_ES, nil
	case "unspecified":
		return Language_Unspecified, nil
	}
	return "", fmt.Errorf("unknown Language: %q", s)
}

type PACTAVersionID string
type PACTAVersion struct {
	ID          PACTAVersionID
	Name        string
	Description string
	Digest      string
	CreatedAt   time.Time
	IsDefault   bool
}

func (o *PACTAVersion) Clone() *PACTAVersion {
	if o == nil {
		return nil
	}
	return &PACTAVersion{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		Digest:      o.Digest,
		CreatedAt:   o.CreatedAt,
		IsDefault:   o.IsDefault,
	}
}

type UserID string
type User struct {
	ID                UserID
	AuthnMechanism    AuthnMechanism
	AuthnID           string
	EnteredEmail      string
	CanonicalEmail    string
	Admin             bool
	SuperAdmin        bool
	Name              string
	PreferredLanguage Language
	CountryCode       string
	CreatedAt         time.Time
}

func (o *User) Clone() *User {
	if o == nil {
		return nil
	}
	return &User{
		ID:                o.ID,
		AuthnMechanism:    o.AuthnMechanism,
		AuthnID:           o.AuthnID,
		EnteredEmail:      o.EnteredEmail,
		CanonicalEmail:    o.CanonicalEmail,
		Admin:             o.Admin,
		SuperAdmin:        o.SuperAdmin,
		Name:              o.Name,
		PreferredLanguage: o.PreferredLanguage,
		CountryCode:       o.CountryCode,
		CreatedAt:         o.CreatedAt,
	}
}

type InitiativeID string
type Initiative struct {
	ID                       InitiativeID
	Name                     string
	Affiliation              string
	PublicDescription        string
	InternalDescription      string
	RequiresInvitationToJoin bool
	IsAcceptingNewMembers    bool
	IsAcceptingNewPortfolios bool
	PACTAVersion             *PACTAVersion
	Language                 Language
	CreatedAt                time.Time
	UserRelationships        []*InitiativeUserRelationship
	PortfolioRelationships   []*PortfolioInitiativeMembership
	Invitations              []*InitiativeInvitation
}

func (o *Initiative) Clone() *Initiative {
	if o == nil {
		return nil
	}
	return &Initiative{
		ID:                       o.ID,
		Name:                     o.Name,
		Affiliation:              o.Affiliation,
		PublicDescription:        o.PublicDescription,
		InternalDescription:      o.InternalDescription,
		RequiresInvitationToJoin: o.RequiresInvitationToJoin,
		IsAcceptingNewMembers:    o.IsAcceptingNewMembers,
		IsAcceptingNewPortfolios: o.IsAcceptingNewPortfolios,
		PACTAVersion:             o.PACTAVersion.Clone(),
		Language:                 o.Language,
		CreatedAt:                o.CreatedAt,
		UserRelationships:        cloneAll(o.UserRelationships),
		PortfolioRelationships:   cloneAll(o.PortfolioRelationships),
		Invitations:              cloneAll(o.Invitations),
	}
}

type InitiativeInvitationID string
type InitiativeInvitation struct {
	ID         InitiativeInvitationID
	Initiative *Initiative
	CreatedAt  time.Time
	UsedAt     time.Time
	UsedBy     *User
}

func (o *InitiativeInvitation) Clone() *InitiativeInvitation {
	if o == nil {
		return nil
	}
	return &InitiativeInvitation{
		ID:         o.ID,
		Initiative: o.Initiative.Clone(),
		CreatedAt:  o.CreatedAt,
		UsedAt:     o.UsedAt,
		UsedBy:     o.UsedBy.Clone(),
	}
}

type InitiativeUserRelationship struct {
	Initiative *Initiative
	User       *User
	Manager    bool
	Member     bool
	UpdatedAt  time.Time
}

func (o *InitiativeUserRelationship) Clone() *InitiativeUserRelationship {
	if o == nil {
		return nil
	}
	return &InitiativeUserRelationship{
		Initiative: o.Initiative.Clone(),
		User:       o.User.Clone(),
		Manager:    o.Manager,
		Member:     o.Member,
		UpdatedAt:  o.UpdatedAt,
	}
}

type FileType string

const (
	FileType_CSV  = "csv"
	FileType_YAML = "yaml"
	FileType_ZIP  = "zip"
	FileType_HTML = "html"
	FileType_JSON = "json"
)

var FileTypeValues = []FileType{
	FileType_CSV,
	FileType_YAML,
	FileType_ZIP,
	FileType_JSON,
	FileType_HTML,
	FileType_JSON,
}

func ParseFileType(s string) (FileType, error) {
	ss := strings.TrimSpace(strings.ToLower(s))
	if strings.HasPrefix(ss, ".") {
		ss = ss[1:]
	}
	switch ss {
	case "csv":
		return FileType_CSV, nil
	case "yaml":
		return FileType_YAML, nil
	case "zip":
		return FileType_ZIP, nil
	case "html":
		return FileType_HTML, nil
	case "json":
		return FileType_JSON, nil
	}
	return "", fmt.Errorf("unknown pacta.FileType: %q", s)
}

type BlobURI string
type BlobID string
type Blob struct {
	ID        BlobID
	BlobURI   BlobURI
	FileType  FileType
	FileName  string
	CreatedAt time.Time
}

func (o *Blob) Clone() *Blob {
	if o == nil {
		return nil
	}
	return &Blob{
		ID:        o.ID,
		BlobURI:   o.BlobURI,
		FileType:  o.FileType,
		FileName:  o.FileName,
		CreatedAt: o.CreatedAt,
	}
}

type OwnerID string
type Owner struct {
	ID         OwnerID
	User       *User
	Initiative *Initiative
}

func (o *Owner) Clone() *Owner {
	if o == nil {
		return nil
	}
	return &Owner{
		ID:         o.ID,
		User:       o.User.Clone(),
		Initiative: o.Initiative.Clone(),
	}
}

type FailureCode string

const (
	FailureCode_Unknown FailureCode = "UNKNOWN"
)

var FailureCodeValues = []FailureCode{
	FailureCode_Unknown,
}

func ParseFailureCode(s string) (FailureCode, error) {
	switch s {
	case "UNKNOWN":
		return FailureCode_Unknown, nil
	}
	return "", fmt.Errorf("unknown FailureCode: %q", s)
}

type HoldingsDate struct {
	Time time.Time
}

func (o *HoldingsDate) Clone() *HoldingsDate {
	if o == nil {
		return nil
	}
	return &HoldingsDate{
		Time: o.Time,
	}
}

type IncompleteUploadID string
type IncompleteUpload struct {
	ID                IncompleteUploadID
	Name              string
	Description       string
	CreatedAt         time.Time
	HoldingsDate      *HoldingsDate
	RanAt             time.Time
	CompletedAt       time.Time
	FailureCode       FailureCode
	FailureMessage    string
	AdminDebugEnabled bool
	Owner             *Owner
	Blob              *Blob
}

func (o *IncompleteUpload) Clone() *IncompleteUpload {
	if o == nil {
		return nil
	}
	return &IncompleteUpload{
		ID:                o.ID,
		Name:              o.Name,
		Description:       o.Description,
		CreatedAt:         o.CreatedAt,
		HoldingsDate:      o.HoldingsDate.Clone(),
		RanAt:             o.RanAt,
		CompletedAt:       o.CompletedAt,
		FailureCode:       o.FailureCode,
		FailureMessage:    o.FailureMessage,
		AdminDebugEnabled: o.AdminDebugEnabled,
		Owner:             o.Owner.Clone(),
		Blob:              o.Blob.Clone(),
	}
}

type PortfolioID string
type Portfolio struct {
	ID                PortfolioID
	Name              string
	Description       string
	CreatedAt         time.Time
	HoldingsDate      *HoldingsDate
	Owner             *Owner
	Blob              *Blob
	AdminDebugEnabled bool
	NumberOfRows      int
}

func (o *Portfolio) Clone() *Portfolio {
	if o == nil {
		return nil
	}
	return &Portfolio{
		ID:                o.ID,
		Name:              o.Name,
		Description:       o.Description,
		CreatedAt:         o.CreatedAt,
		HoldingsDate:      o.HoldingsDate.Clone(),
		Owner:             o.Owner.Clone(),
		Blob:              o.Blob.Clone(),
		AdminDebugEnabled: o.AdminDebugEnabled,
		NumberOfRows:      o.NumberOfRows,
	}
}

type PortfolioGroupID string
type PortfolioGroup struct {
	ID          PortfolioGroupID
	Owner       *Owner
	Name        string
	Description string
	CreatedAt   time.Time
	Members     []*PortfolioGroupMembership
}

func (o *PortfolioGroup) Clone() *PortfolioGroup {
	if o == nil {
		return nil
	}
	return &PortfolioGroup{
		ID:          o.ID,
		Owner:       o.Owner.Clone(),
		Name:        o.Name,
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		Members:     cloneAll(o.Members),
	}
}

type PortfolioGroupMembership struct {
	PortfolioGroup *PortfolioGroup
	Portfolio      *Portfolio
	CreatedAt      time.Time
}

func (o *PortfolioGroupMembership) Clone() *PortfolioGroupMembership {
	if o == nil {
		return nil
	}
	return &PortfolioGroupMembership{
		PortfolioGroup: o.PortfolioGroup.Clone(),
		Portfolio:      o.Portfolio.Clone(),
		CreatedAt:      o.CreatedAt,
	}
}

type PortfolioInitiativeMembership struct {
	Portfolio  *Portfolio
	Initiative *Initiative
	CreatedAt  time.Time
	AddedBy    *User
}

func (o *PortfolioInitiativeMembership) Clone() *PortfolioInitiativeMembership {
	if o == nil {
		return nil
	}
	return &PortfolioInitiativeMembership{
		Portfolio:  o.Portfolio.Clone(),
		Initiative: o.Initiative.Clone(),
		CreatedAt:  o.CreatedAt,
		AddedBy:    o.AddedBy.Clone(),
	}
}

type PortfolioSnapshotID string
type PortfolioSnapshot struct {
	ID               PortfolioSnapshotID
	PortfolioID      PortfolioID
	PortfolioGroupID PortfolioGroupID
	InitiatiativeID  InitiativeID
	PortfolioIDs     []PortfolioID
}

func (o *PortfolioSnapshot) Clone() *PortfolioSnapshot {
	if o == nil {
		return nil
	}
	pids := make([]PortfolioID, len(o.PortfolioIDs))
	copy(pids, o.PortfolioIDs)
	return &PortfolioSnapshot{
		ID:               o.ID,
		PortfolioID:      o.PortfolioID,
		PortfolioGroupID: o.PortfolioGroupID,
		InitiatiativeID:  o.InitiatiativeID,
		PortfolioIDs:     pids,
	}
}

type AnalysisType string

const (
	AnalysisType_Audit  AnalysisType = "audit"
	AnalysisType_Report AnalysisType = "report"
)

var AnalysisTypeValues = []AnalysisType{
	AnalysisType_Audit,
	AnalysisType_Report,
}

func ParseAnalysisType(s string) (AnalysisType, error) {
	switch s {
	case "audit":
		return AnalysisType_Audit, nil
	case "report":
		return AnalysisType_Report, nil
	}
	return "", fmt.Errorf("unknown AnalysisType: %q", s)
}

type AnalysisID string
type Analysis struct {
	ID                AnalysisID
	AnalysisType      AnalysisType
	Owner             *Owner
	PACTAVersion      *PACTAVersion
	PortfolioSnapshot *PortfolioSnapshot
	Name              string
	Description       string
	CreatedAt         time.Time
	RanAt             time.Time
	CompletedAt       time.Time
	FailureCode       FailureCode
	FailureMessage    string
	Artifacts         []*AnalysisArtifact
}

func (o *Analysis) Clone() *Analysis {
	if o == nil {
		return nil
	}
	return &Analysis{
		ID:                o.ID,
		AnalysisType:      o.AnalysisType,
		Owner:             o.Owner.Clone(),
		PACTAVersion:      o.PACTAVersion.Clone(),
		PortfolioSnapshot: o.PortfolioSnapshot.Clone(),
		Name:              o.Name,
		Description:       o.Description,
		CreatedAt:         o.CreatedAt,
		RanAt:             o.RanAt,
		CompletedAt:       o.CompletedAt,
		FailureCode:       o.FailureCode,
		FailureMessage:    o.FailureMessage,
		Artifacts:         cloneAll(o.Artifacts),
	}
}

type AnalysisArtifactID string
type AnalysisArtifact struct {
	ID                AnalysisArtifactID
	AnalysisID        AnalysisID
	Blob              *Blob
	AdminDebugEnabled bool
	SharedToPublic    bool
}

func (o *AnalysisArtifact) Clone() *AnalysisArtifact {
	if o == nil {
		return nil
	}
	return &AnalysisArtifact{
		ID:                o.ID,
		AnalysisID:        o.AnalysisID,
		Blob:              o.Blob.Clone(),
		AdminDebugEnabled: o.AdminDebugEnabled,
		SharedToPublic:    o.SharedToPublic,
	}
}

type AuditLogAction string

const (
	AuditLogAction_Create            AuditLogAction = "CREATE"
	AuditLogAction_Update            AuditLogAction = "UPDATE"
	AuditLogAction_Delete            AuditLogAction = "DELETE"
	AuditLogAction_AddTo             AuditLogAction = "ADD_TO"
	AuditLogAction_RemoveFrom        AuditLogAction = "REMOVE_FROM"
	AuditLogAction_EnableAdminDebug  AuditLogAction = "ENABLE_ADMIN_DEBUG"
	AuditLogAction_DisableAdminDebug AuditLogAction = "DISABLE_ADMIN_DEBUG"
	AuditLogAction_Download          AuditLogAction = "DOWNLOAD"
	AuditLogAction_EnableSharing     AuditLogAction = "ENABLE_SHARING"
	AuditLogAction_DisableSharing    AuditLogAction = "DISABLE_SHARING"
)

var AuditLogActionValues = []AuditLogAction{
	AuditLogAction_Create,
	AuditLogAction_Update,
	AuditLogAction_Delete,
	AuditLogAction_AddTo,
	AuditLogAction_RemoveFrom,
	AuditLogAction_EnableAdminDebug,
	AuditLogAction_DisableAdminDebug,
	AuditLogAction_Download,
	AuditLogAction_EnableSharing,
	AuditLogAction_DisableSharing,
}

func ParseAuditLogAction(s string) (AuditLogAction, error) {
	switch s {
	case "CREATE":
		return AuditLogAction_Create, nil
	case "UPDATE":
		return AuditLogAction_Update, nil
	case "DELETE":
		return AuditLogAction_Delete, nil
	case "ADD_TO":
		return AuditLogAction_AddTo, nil
	case "REMOVE_FROM":
		return AuditLogAction_RemoveFrom, nil
	case "ENABLE_ADMIN_DEBUG":
		return AuditLogAction_EnableAdminDebug, nil
	case "DISABLE_ADMIN_DEBUG":
		return AuditLogAction_DisableAdminDebug, nil
	case "DOWNLOAD":
		return AuditLogAction_Download, nil
	case "ENABLE_SHARING":
		return AuditLogAction_EnableSharing, nil
	case "DISABLE_SHARING":
		return AuditLogAction_DisableSharing, nil
	}
	return "", fmt.Errorf("unknown AuditLogAction: %q", s)
}

type AuditLogActorType string

const (
	AuditLogActorType_User       AuditLogActorType = "USER"
	AuditLogActorType_Admin      AuditLogActorType = "ADMIN"
	AuditLogActorType_SuperAdmin AuditLogActorType = "SUPER_ADMIN"
	AuditLogActorType_System     AuditLogActorType = "SYSTEM"
)

var AuditLogActorTypeValues = []AuditLogActorType{
	AuditLogActorType_User,
	AuditLogActorType_Admin,
	AuditLogActorType_SuperAdmin,
	AuditLogActorType_System,
}

func ParseAuditLogActorType(s string) (AuditLogActorType, error) {
	switch s {
	case "USER":
		return AuditLogActorType_User, nil
	case "ADMIN":
		return AuditLogActorType_Admin, nil
	case "SUPER_ADMIN":
		return AuditLogActorType_SuperAdmin, nil
	case "SYSTEM":
		return AuditLogActorType_System, nil
	}
	return "", fmt.Errorf("unknown AuditLogActorType: %q", s)
}

type AuditLogTargetType string

const (
	AuditLogTargetType_User             AuditLogTargetType = "USER"
	AuditLogTargetType_Portfolio        AuditLogTargetType = "PORTFOLIO"
	AuditLogTargetType_IncompleteUpload AuditLogTargetType = "INCOMPLETE_UPLOAD"
	AuditLogTargetType_PortfolioGroup   AuditLogTargetType = "PORTFOLIO_GROUP"
	AuditLogTargetType_Initiative       AuditLogTargetType = "INITIATIVE"
	AuditLogTargetType_PACTAVersion     AuditLogTargetType = "PACTA_VERSION"
	AuditLogTargetType_Analysis         AuditLogTargetType = "ANALYSIS"
)

var AuditLogTargetTypeValues = []AuditLogTargetType{
	AuditLogTargetType_User,
	AuditLogTargetType_Portfolio,
	AuditLogTargetType_IncompleteUpload,
	AuditLogTargetType_PortfolioGroup,
	AuditLogTargetType_Initiative,
	AuditLogTargetType_PACTAVersion,
	AuditLogTargetType_Analysis,
}

func ParseAuditLogTargetType(s string) (AuditLogTargetType, error) {
	switch s {
	case "USER":
		return AuditLogTargetType_User, nil
	case "PORTFOLIO":
		return AuditLogTargetType_Portfolio, nil
	case "INCOMPLETE_UPLOAD":
		return AuditLogTargetType_IncompleteUpload, nil
	case "PORTFOLIO_GROUP":
		return AuditLogTargetType_PortfolioGroup, nil
	case "INITIATIVE":
		return AuditLogTargetType_Initiative, nil
	case "PACTA_VERSION":
		return AuditLogTargetType_PACTAVersion, nil
	case "ANALYSIS":
		return AuditLogTargetType_Analysis, nil
	}
	return "", fmt.Errorf("unknown AuditLogTargetType: %q", s)
}

type AuditLogID string
type AuditLog struct {
	ID                   AuditLogID
	CreatedAt            time.Time
	ActorType            AuditLogActorType
	ActorID              string
	ActorOwner           *Owner
	Action               AuditLogAction
	PrimaryTargetType    AuditLogTargetType
	PrimaryTargetID      string
	PrimaryTargetOwner   *Owner
	SecondaryTargetType  AuditLogTargetType
	SecondaryTargetID    string
	SecondaryTargetOwner *Owner
}

func (o *AuditLog) Clone() *AuditLog {
	if o == nil {
		return nil
	}
	return &AuditLog{
		ID:                   o.ID,
		CreatedAt:            o.CreatedAt,
		ActorType:            o.ActorType,
		ActorID:              o.ActorID,
		ActorOwner:           o.ActorOwner.Clone(),
		Action:               o.Action,
		PrimaryTargetType:    o.PrimaryTargetType,
		PrimaryTargetID:      o.PrimaryTargetID,
		PrimaryTargetOwner:   o.PrimaryTargetOwner.Clone(),
		SecondaryTargetType:  o.SecondaryTargetType,
		SecondaryTargetID:    o.SecondaryTargetID,
		SecondaryTargetOwner: o.SecondaryTargetOwner.Clone(),
	}
}

type cloneable[T any] interface {
	comparable
	Clone() T
}

func cloneAll[T cloneable[T]](in []T) []T {
	out := make([]T, len(in))
	for i, t := range in {
		out[i] = t.Clone()
	}
	return out
}
