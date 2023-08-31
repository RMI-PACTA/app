package pacta

import "testing"

func TestParseAuthMechanism(t *testing.T) {
	testParseEnum(t, AuthnMechanismValues, ParseAuthnMechanism)
}

func TestParseLanguage(t *testing.T) {
	testParseEnum(t, LanguageValues, ParseLanguage)
}

func TestParseFileType(t *testing.T) {
	testParseEnum(t, FileTypeValues, ParseFileType)
}

func TestParseFailureCode(t *testing.T) {
	testParseEnum(t, FailureCodeValues, ParseFailureCode)
}

func TestParseAnalysisType(t *testing.T) {
	testParseEnum(t, AnalysisTypeValues, ParseAnalysisType)
}

func TestParseAuditLogAction(t *testing.T) {
	testParseEnum(t, AuditLogActionValues, ParseAuditLogAction)
}

func TestParseAuditLogActorType(t *testing.T) {
	testParseEnum(t, AuditLogActorTypeValues, ParseAuditLogActorType)
}

func TestParseAuditLogTargetType(t *testing.T) {
	testParseEnum(t, AuditLogTargetTypeValues, ParseAuditLogTargetType)
}

func testParseEnum[E ~string](t *testing.T, es []E, fn func(string) (E, error)) {
	t.Helper()
	for _, e := range es {
		s := string(e)
		e2, err := fn(s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if e != e2 {
			t.Errorf("expected %v, got %v", e, e2)
		}
	}
	e, err := fn("invalid")
	if err == nil {
		t.Errorf("expected error, got %v", e)
	}
}
