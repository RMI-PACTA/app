package conv

import (
	"fmt"

	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

func LanguageToOAPI(l pacta.Language) (api.Language, error) {
	switch l {
	case pacta.Language_DE:
		return api.LanguageDE, nil
	case pacta.Language_ES:
		return api.LanguageES, nil
	case pacta.Language_EN:
		return api.LanguageEN, nil
	case pacta.Language_FR:
		return api.LanguageFR, nil
	default:
		return "", fmt.Errorf("unknown language: %q", l)
	}
}

func InitiativeToOAPI(i *pacta.Initiative) (*api.Initiative, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeToOAPI: can't convert nil pointer")
	}
	pims, err := convAll(i.PortfolioInitiativeMemberships, portfolioInitiativeMembershipToOAPIPortfolio)
	if err != nil {
		return nil, oapierr.Internal("initiativeToOAPI: portfolioInitiativeMembershipToOAPIInitiative failed", zap.Error(err))
	}
	lang, err := LanguageToOAPI(i.Language)
	if err != nil {
		return nil, oapierr.Internal("initiativeToOAPI: languageToOAPI failed", zap.Error(err))
	}
	return &api.Initiative{
		Affiliation:                    i.Affiliation,
		CreatedAt:                      i.CreatedAt,
		Id:                             string(i.ID),
		InternalDescription:            i.InternalDescription,
		IsAcceptingNewMembers:          i.IsAcceptingNewMembers,
		IsAcceptingNewPortfolios:       i.IsAcceptingNewPortfolios,
		Language:                       lang,
		Name:                           i.Name,
		PactaVersion:                   strPtr(i.PACTAVersion.ID),
		PublicDescription:              i.PublicDescription,
		RequiresInvitationToJoin:       i.RequiresInvitationToJoin,
		PortfolioInitiativeMemberships: pims,
	}, nil
}

func portfolioInitiativeMembershipToOAPIPortfolio(in *pacta.PortfolioInitiativeMembership) (api.PortfolioInitiativeMembershipPortfolio, error) {
	var zero api.PortfolioInitiativeMembershipPortfolio
	out := &api.PortfolioInitiativeMembershipPortfolio{
		CreatedAt: in.CreatedAt,
	}
	if in.AddedBy != nil && in.AddedBy.ID == "" {
		out.AddedByUserId = strPtr(in.AddedBy.ID)
	}
	p, err := PortfolioToOAPI(in.Portfolio)
	if err != nil {
		return zero, oapierr.Internal("portfolioInitiativeMembershipToOAPI: portfolioToOAPI failed", zap.Error(err))
	}
	out.Portfolio = *p
	return zero, nil
}

func portfolioInitiativeMembershipToOAPIInitiative(in *pacta.PortfolioInitiativeMembership) (api.PortfolioInitiativeMembershipInitiative, error) {
	var zero api.PortfolioInitiativeMembershipInitiative
	out := api.PortfolioInitiativeMembershipInitiative{
		CreatedAt: in.CreatedAt,
	}
	if in.AddedBy != nil && in.AddedBy.ID == "" {
		out.AddedByUserId = strPtr(in.AddedBy.ID)
	}
	if in.Initiative != nil && in.Initiative.PACTAVersion != nil {
		i, err := InitiativeToOAPI(in.Initiative)
		if err != nil {
			return zero, oapierr.Internal("portfolioInitiativeMembershipToOAPI: initiativeToOAPI failed", zap.Error(err))
		}
		out.Initiative = *i
	}
	return out, nil
}

func UserToOAPI(user *pacta.User) (*api.User, error) {
	if user == nil {
		return nil, oapierr.Internal("userToOAPI: can't convert nil pointer")
	}
	var lang api.Language
	if user.PreferredLanguage != "" {
		l, err := LanguageToOAPI(user.PreferredLanguage)
		if err != nil {
			return nil, oapierr.Internal("userToOAPI: languageToOAPI failed", zap.Error(err))
		}
		lang = l
	}
	return &api.User{
		Admin:             user.Admin,
		CanonicalEmail:    &user.CanonicalEmail,
		EnteredEmail:      user.EnteredEmail,
		Id:                string(user.ID),
		Name:              user.Name,
		PreferredLanguage: lang,
		SuperAdmin:        user.SuperAdmin,
	}, nil
}

func PactaVersionToOAPI(pv *pacta.PACTAVersion) (*api.PactaVersion, error) {
	if pv == nil {
		return nil, oapierr.Internal("pactaVersionToOAPI: can't convert nil pointer")
	}
	return &api.PactaVersion{
		CreatedAt:   pv.CreatedAt,
		Description: pv.Description,
		Digest:      pv.Digest,
		Id:          string(pv.ID),
		IsDefault:   pv.IsDefault,
		Name:        pv.Name,
	}, nil
}

func InitiativeInvitationToOAPI(i *pacta.InitiativeInvitation) (*api.InitiativeInvitation, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeToOAPI: can't convert nil pointer")
	}
	var usedAt *string
	if !i.UsedAt.IsZero() {
		usedAt = ptr(i.UsedAt.String())
	}
	var usedBy *string
	if i.UsedBy != nil {
		usedBy = strPtr(i.UsedBy.ID)
	}
	return &api.InitiativeInvitation{
		CreatedAt:    i.CreatedAt,
		Id:           string(i.ID),
		InitiativeId: string(i.Initiative.ID),
		UsedAt:       usedAt,
		UsedByUserId: usedBy,
	}, nil
}

func InitiativeUserRelationshipToOAPI(i *pacta.InitiativeUserRelationship) (*api.InitiativeUserRelationship, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeUserRelationshipToOAPI: can't convert nil pointer")
	}
	if i.User == nil {
		return nil, oapierr.Internal("initiativeUserRelationshipToOAPI: can't convert nil user")
	}
	if i.Initiative == nil {
		return nil, oapierr.Internal("initiativeUserRelationshipToOAPI: can't convert nil initiative")
	}
	return &api.InitiativeUserRelationship{
		UpdatedAt:    i.UpdatedAt,
		InitiativeId: string(i.Initiative.ID),
		UserId:       string(i.User.ID),
		Manager:      i.Manager,
		Member:       i.Member,
	}, nil
}

func HoldingsDateToOAPI(hd *pacta.HoldingsDate) (*api.HoldingsDate, error) {
	if hd == nil {
		return nil, nil
	}
	return &api.HoldingsDate{
		Time: hd.Time,
	}, nil
}

func IncompleteUploadsToOAPI(ius []*pacta.IncompleteUpload) ([]*api.IncompleteUpload, error) {
	return convAll(ius, IncompleteUploadToOAPI)
}

func FailureCodeToOAPI(f pacta.FailureCode) (*api.FailureCode, error) {
	switch f {
	case "":
		return nil, nil
	case pacta.FailureCode_Unknown:
		return ptr(api.FailureCodeUNKNOWN), nil
	}
	return nil, fmt.Errorf("unknown failure code: %q", f)
}

func IncompleteUploadToOAPI(iu *pacta.IncompleteUpload) (*api.IncompleteUpload, error) {
	if iu == nil {
		return nil, oapierr.Internal("incompleteUploadToOAPI: can't convert nil pointer")
	}
	hd, err := HoldingsDateToOAPI(iu.HoldingsDate)
	if err != nil {
		return nil, oapierr.Internal("incompleteUploadToOAPI: holdingsDateToOAPI failed", zap.Error(err))
	}
	fc, err := FailureCodeToOAPI(iu.FailureCode)
	if err != nil {
		return nil, oapierr.Internal("incompleteUploadToOAPI: failureCodeToOAPI failed", zap.Error(err))
	}
	return &api.IncompleteUpload{
		Id:                string(iu.ID),
		Name:              iu.Name,
		Description:       iu.Description,
		HoldingsDate:      hd,
		CreatedAt:         iu.CreatedAt,
		RanAt:             timeToNilable(iu.RanAt),
		CompletedAt:       timeToNilable(iu.CompletedAt),
		FailureCode:       fc,
		FailureMessage:    stringToNilable(iu.FailureMessage),
		AdminDebugEnabled: iu.AdminDebugEnabled,
	}, nil
}

func PortfoliosToOAPI(ius []*pacta.Portfolio) ([]*api.Portfolio, error) {
	return convAll(ius, PortfolioToOAPI)
}

func PortfolioToOAPI(p *pacta.Portfolio) (*api.Portfolio, error) {
	if p == nil {
		return nil, oapierr.Internal("portfolioToOAPI: can't convert nil pointer")
	}
	hd, err := HoldingsDateToOAPI(p.HoldingsDate)
	if err != nil {
		return nil, oapierr.Internal("portfolioToOAPI: holdingsDateToOAPI failed", zap.Error(err))
	}
	portfolioGroupMemberships := []api.PortfolioGroupMembershipPortfolioGroup{}
	for _, m := range p.PortfolioGroupMemberships {
		pg, err := PortfolioGroupToOAPI(m.PortfolioGroup)
		if err != nil {
			return nil, oapierr.Internal("portfolioToOAPI: portfolioGroupToOAPI failed", zap.Error(err))
		}
		portfolioGroupMemberships = append(portfolioGroupMemberships, api.PortfolioGroupMembershipPortfolioGroup{
			CreatedAt:      m.CreatedAt,
			PortfolioGroup: *pg,
		})
	}
	pims, err := convAll(p.PortfolioInitiativeMemberships, portfolioInitiativeMembershipToOAPIInitiative)
	if err != nil {
		return nil, oapierr.Internal("initiativeToOAPI: portfolioInitiativeMembershipToOAPIInitiative failed", zap.Error(err))
	}
	return &api.Portfolio{
		Id:                string(p.ID),
		Name:              p.Name,
		Description:       p.Description,
		HoldingsDate:      hd,
		CreatedAt:         p.CreatedAt,
		NumberOfRows:      p.NumberOfRows,
		AdminDebugEnabled: p.AdminDebugEnabled,
		Groups:            &portfolioGroupMemberships,
		Initiatives:       &pims,
	}, nil
}

func PortfolioGroupToOAPI(pg *pacta.PortfolioGroup) (*api.PortfolioGroup, error) {
	if pg == nil {
		return nil, oapierr.Internal("portfolioGroupToOAPI: can't convert nil pointer")
	}
	members := []api.PortfolioGroupMembershipPortfolio{}
	for _, m := range pg.PortfolioGroupMemberships {
		portfolio, err := PortfolioToOAPI(m.Portfolio)
		if err != nil {
			return nil, oapierr.Internal("portfolioGroupToOAPI: portfolioToOAPI failed", zap.Error(err))
		}
		members = append(members, api.PortfolioGroupMembershipPortfolio{
			CreatedAt: m.CreatedAt,
			Portfolio: *portfolio,
		})
	}
	return &api.PortfolioGroup{
		Id:          string(pg.ID),
		Name:        pg.Name,
		Description: pg.Description,
		CreatedAt:   pg.CreatedAt,
		Members:     &members,
	}, nil
}

func PortfolioGroupsToOAPI(pgs []*pacta.PortfolioGroup) ([]*api.PortfolioGroup, error) {
	return convAll(pgs, PortfolioGroupToOAPI)
}

func FileTypeToOAPI(ft pacta.FileType) (api.FileType, error) {
	switch ft {
	case pacta.FileType_CSV:
		return api.FileTypeCSV, nil
	case pacta.FileType_YAML:
		return api.FileTypeYAML, nil
	case pacta.FileType_ZIP:
		return api.FileTypeZIP, nil
	case pacta.FileType_JSON:
		return api.FileTypeJSON, nil
	case pacta.FileType_HTML:
		return api.FileTypeHTML, nil
	}
	return "", fmt.Errorf("unknown file type: %q", ft)
}

func BlobToOAPI(b *pacta.Blob) (*api.Blob, error) {
	if b == nil {
		return nil, oapierr.Internal("blobToOAPI: can't convert nil pointer")
	}
	ft, err := FileTypeToOAPI(b.FileType)
	if err != nil {
		return nil, oapierr.Internal("blobToOAPI: fileTypeToOAPI failed", zap.Error(err))
	}
	return &api.Blob{
		Id:        string(b.ID),
		FileName:  b.FileName,
		FileType:  ft,
		CreatedAt: b.CreatedAt,
	}, nil
}

func PortfolioSnapshotToOAPI(ps *pacta.PortfolioSnapshot) (*api.PortfolioSnapshot, error) {
	if ps == nil {
		return nil, oapierr.Internal("portfolioSnapshotToOAPI: can't convert nil pointer")
	}
	portfolioIds := make([]string, len(ps.PortfolioIDs))
	for i, pid := range ps.PortfolioIDs {
		portfolioIds[i] = string(pid)
	}
	result := &api.PortfolioSnapshot{
		Id:           string(ps.ID),
		PortfolioIds: portfolioIds,
	}
	if ps.Portfolio != nil {
		portfolio, err := PortfolioToOAPI(ps.Portfolio)
		if err != nil {
			return nil, oapierr.Internal("portfolioSnapshotToOAPI: portfolioToOAPI failed", zap.Error(err))
		}
		result.Portfolio = portfolio
	}
	if ps.PortfolioGroup != nil {
		portfolioGroup, err := PortfolioGroupToOAPI(ps.PortfolioGroup)
		if err != nil {
			return nil, oapierr.Internal("portfolioSnapshotToOAPI: portfolioGroupToOAPI failed", zap.Error(err))
		}
		result.PortfolioGroup = portfolioGroup
	}
	if ps.Initiatiative != nil {
		initiative, err := InitiativeToOAPI(ps.Initiatiative)
		if err != nil {
			return nil, oapierr.Internal("portfolioSnapshotToOAPI: initiativeToOAPI failed", zap.Error(err))
		}
		result.Initiative = initiative
	}
	return result, nil
}

func AnalysisArtifactToOAPI(aa *pacta.AnalysisArtifact) (*api.AnalysisArtifact, error) {
	if aa == nil {
		return nil, oapierr.Internal("analysisArtifactToOAPI: can't convert nil pointer")
	}
	blob, err := BlobToOAPI(aa.Blob)
	if err != nil {
		return nil, oapierr.Internal("analysisArtifactToOAPI: blobToOAPI failed", zap.Error(err))
	}
	return &api.AnalysisArtifact{
		Id:                string(aa.ID),
		AdminDebugEnabled: aa.AdminDebugEnabled,
		SharedToPublic:    aa.AdminDebugEnabled,
		Blob:              *blob,
	}, nil
}

func AnalysisTypeToOAPI(at pacta.AnalysisType) (api.AnalysisType, error) {
	switch at {
	case pacta.AnalysisType_Audit:
		return api.AnalysisTypeAUDIT, nil
	case pacta.AnalysisType_Report:
		return api.AnalysisTypeREPORT, nil
	}
	return "", fmt.Errorf("unknown analysis type: %q", at)
}

func AnalysisToOAPI(a *pacta.Analysis) (*api.Analysis, error) {
	if a == nil {
		return nil, oapierr.Internal("analysisToOAPI: can't convert nil pointer")
	}
	aas, err := convAll(a.Artifacts, AnalysisArtifactToOAPI)
	if err != nil {
		return nil, oapierr.Internal("analysisToOAPI: analysisArtifactsToOAPI failed", zap.Error(err))
	}
	snapshot, err := PortfolioSnapshotToOAPI(a.PortfolioSnapshot)
	if err != nil {
		return nil, oapierr.Internal("analysisToOAPI: portfolioSnapshotToOAPI failed", zap.Error(err))
	}
	var fc *api.FailureCode
	if a.FailureCode != "" {
		fc = ptr(api.FailureCode(a.FailureCode))
	}
	var fm *string
	if a.FailureMessage != "" {
		fm = ptr(a.FailureMessage)
	}
	at, err := AnalysisTypeToOAPI(a.AnalysisType)
	if err != nil {
		return nil, oapierr.Internal("analysisToOAPI: analysisTypeToOAPI failed", zap.Error(err))
	}
	return &api.Analysis{
		Id:                string(a.ID),
		AnalysisType:      at,
		PactaVersion:      string(a.PACTAVersion.ID),
		PortfolioSnapshot: *snapshot,
		Name:              a.Name,
		Description:       a.Description,
		CreatedAt:         a.CreatedAt,
		RanAt:             timeToNilable(a.RanAt),
		CompletedAt:       timeToNilable(a.CompletedAt),
		FailureCode:       fc,
		FailureMessage:    fm,
		Artifacts:         dereferenceAll(aas),
	}, nil
}

func AnalysesToOAPI(as []*pacta.Analysis) ([]*api.Analysis, error) {
	return convAll(as, AnalysisToOAPI)
}

func auditLogActorTypeToOAPI(i pacta.AuditLogActorType) (api.AuditLogActorType, error) {
	switch i {
	case pacta.AuditLogActorType_Public:
		return api.AuditLogActorTypePublic, nil
	case pacta.AuditLogActorType_Owner:
		return api.AuditLogActorTypeOwner, nil
	case pacta.AuditLogActorType_Admin:
		return api.AuditLogActorTypeAdmin, nil
	case pacta.AuditLogActorType_SuperAdmin:
		return api.AuditLogActorTypeSuperAdmin, nil
	case pacta.AuditLogActorType_System:
		return api.AuditLogActorTypeSystem, nil
	}
	return "", oapierr.Internal(fmt.Sprintf("auditLogActorTypeToOAPI: unknown actor type: %q", i))
}

func auditLogActionToOAPI(i pacta.AuditLogAction) (api.AuditLogAction, error) {
	switch i {
	case pacta.AuditLogAction_Create:
		return api.AuditLogActionCreate, nil
	case pacta.AuditLogAction_Update:
		return api.AuditLogActionUpdate, nil
	case pacta.AuditLogAction_Delete:
		return api.AuditLogActionDelete, nil
	case pacta.AuditLogAction_AddTo:
		return api.AuditLogActionAddTo, nil
	case pacta.AuditLogAction_RemoveFrom:
		return api.AuditLogActionRemoveFrom, nil
	case pacta.AuditLogAction_EnableAdminDebug:
		return api.AuditLogActionEnableAdminDebug, nil
	case pacta.AuditLogAction_DisableAdminDebug:
		return api.AuditLogActionDisableAdminDebug, nil
	case pacta.AuditLogAction_Download:
		return api.AuditLogActionDownload, nil
	case pacta.AuditLogAction_EnableSharing:
		return api.AuditLogActionEnableSharing, nil
	case pacta.AuditLogAction_DisableSharing:
		return api.AuditLogActionDisableSharing, nil
	case pacta.AuditLogAction_ReadMetadata:
		return api.AuditLogActionReadMetadata, nil
	case pacta.AuditLogAction_TransferOwnership:
		return api.AuditLogActionTransferOwnership, nil
	}
	return "", oapierr.Internal(fmt.Sprintf("auditLogActionToOAPI: unknown action: %q", i))
}

func auditLogTargetTypeToOAPI(i pacta.AuditLogTargetType) (api.AuditLogTargetType, error) {
	switch i {
	case pacta.AuditLogTargetType_User:
		return api.AuditLogTargetTypeUser, nil
	case pacta.AuditLogTargetType_Portfolio:
		return api.AuditLogTargetTypePortfolio, nil
	case pacta.AuditLogTargetType_IncompleteUpload:
		return api.AuditLogTargetTypeIncompleteUpload, nil
	case pacta.AuditLogTargetType_PortfolioGroup:
		return api.AuditLogTargetTypePortfolioGroup, nil
	case pacta.AuditLogTargetType_Initiative:
		return api.AuditLogTargetTypeInitiative, nil
	case pacta.AuditLogTargetType_InitiativeInvitation:
		return api.AuditLogTargetTypeInitiativeInvitation, nil
	case pacta.AuditLogTargetType_PACTAVersion:
		return api.AuditLogTargetTypePactaVersion, nil
	case pacta.AuditLogTargetType_Analysis:
		return api.AuditLogTargetTypeAnalysis, nil
	case pacta.AuditLogTargetType_AnalysisArtifact:
		return api.AuditLogTargetTypeAnalysisArtifact, nil
	}
	return "", oapierr.Internal(fmt.Sprintf("auditLogTargetTypeToOAPI: unknown target type: %q", i))
}

func AuditLogToOAPI(al *pacta.AuditLog) (*api.AuditLog, error) {
	if al == nil {
		return nil, oapierr.Internal("auditLogToOAPI: can't convert nil pointer")
	}
	at, err := auditLogActorTypeToOAPI(al.ActorType)
	if err != nil {
		return nil, oapierr.Internal("auditLogToOAPI: auditLogActorTypeToOAPI failed", zap.Error(err))
	}
	act, err := auditLogActionToOAPI(al.Action)
	if err != nil {
		return nil, oapierr.Internal("auditLogToOAPI: auditLogActionToOAPI failed", zap.Error(err))
	}
	ptt, err := auditLogTargetTypeToOAPI(al.PrimaryTargetType)
	if err != nil {
		return nil, oapierr.Internal("auditLogToOAPI: auditLogTargetTypeToOAPI failed", zap.Error(err))
	}
	var aoi *string
	if al.ActorOwner != nil {
		aoi = stringToNilable(al.ActorOwner.ID)
	}
	var stt *api.AuditLogTargetType
	if al.SecondaryTargetType != "" {
		s, err := auditLogTargetTypeToOAPI(al.SecondaryTargetType)
		if err != nil {
			return nil, oapierr.Internal("auditLogToOAPI: auditLogTargetTypeToOAPI failed", zap.Error(err))
		}
		stt = &s
	}
	var sto *string
	if al.SecondaryTargetOwner != nil {
		s := string(al.SecondaryTargetOwner.ID)
		sto = &s
	}
	var sid *string
	if al.SecondaryTargetID != "" {
		s := string(al.SecondaryTargetID)
		sid = &s
	}
	return &api.AuditLog{
		Id:                   string(al.ID),
		CreatedAt:            al.CreatedAt,
		ActorType:            at,
		ActorId:              stringToNilable(al.ActorID),
		ActorOwnerId:         aoi,
		Action:               act,
		PrimaryTargetType:    ptt,
		PrimaryTargetId:      al.PrimaryTargetID,
		PrimaryTargetOwner:   string(al.PrimaryTargetOwner.ID),
		SecondaryTargetType:  stt,
		SecondaryTargetId:    sid,
		SecondaryTargetOwner: sto,
	}, nil
}

func AuditLogsToOAPI(als []*pacta.AuditLog) ([]*api.AuditLog, error) {
	return convAll(als, AuditLogToOAPI)
}
