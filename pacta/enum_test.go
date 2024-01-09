package pacta

import (
	"path/filepath"
	"testing"
)

func TestParseAuthMechanism(t *testing.T) {
	testParseEnum(t, AuthnMechanismValues, ParseAuthnMechanism)
}

func TestParseLanguage(t *testing.T) {
	testParseEnum(t, LanguageValues, ParseLanguage)
}

func TestParseFileType(t *testing.T) {
	testParseEnum(t, FileTypeValues, ParseFileType)
	otherCases := []string{
		"hello/world.json",
		"hello/world/hithere.JsOn",
		"  hello/world/hithere.json   ",
	}
	for _, c := range otherCases {
		ft, err := ParseFileType(filepath.Ext(c))
		if err != nil {
			t.Errorf("expected successful parse, got %v", err)
		}
		if ft != FileType_JSON {
			t.Errorf("expected JSON, got %v", ft)
		}
	}
}

// need
func TestParseFailureCode(t *testing.T) {
	testParseEnum(t, FailureCodeValues, ParseFailureCode)
}

// need
func TestParseAnalysisType(t *testing.T) {
	testParseEnum(t, AnalysisTypeValues, ParseAnalysisType)
}

// need
func TestParseAuditLogAction(t *testing.T) {
	testParseEnum(t, AuditLogActionValues, ParseAuditLogAction)
}

// need
func TestParseAuditLogActorType(t *testing.T) {
	testParseEnum(t, AuditLogActorTypeValues, ParseAuditLogActorType)
}

// need
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
