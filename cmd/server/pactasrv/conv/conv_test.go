package conv

import (
	"testing"

	"github.com/RMI/pacta/pacta"
)

func TestLanguageRoundTrip(t *testing.T) {
	testEnumConvertability(t, pacta.LanguageValues, LanguageToOAPI, LanguageFromOAPI)
}

func TestAuditLogActorTypeRoundTrip(t *testing.T) {
	testEnumConvertability(t, pacta.AuditLogActorTypeValues, auditLogActorTypeToOAPI, auditLogActorTypeFromOAPI)
}

func TestAuditLogActionRoundTrip(t *testing.T) {
	testEnumConvertability(t, pacta.AuditLogActionValues, auditLogActionToOAPI, auditLogActionFromOAPI)
}

func TestAuditLogTargetTypeRoundTrip(t *testing.T) {
	testEnumConvertability(t, pacta.AuditLogTargetTypeValues, auditLogTargetTypeToOAPI, auditLogTargetTypeFromOAPI)
}

func testEnumConvertability[A comparable, B any](t *testing.T, as []A, aToB func(in A) (B, error), bToA func(in B) (A, error)) {
	for _, a := range as {
		b, err := aToB(a)
		if err != nil {
			t.Fatalf("converting from %T %v: %v", a, a, err)
		}
		a2, err := bToA(b)
		if err != nil {
			t.Fatalf("converting from %T %v: %v", b, b, err)
		}
		if a != a2 {
			t.Errorf("conversion from %T %v to %T %v and back failed, returned %v", a, a, b, b, a2)
		}
	}
}
