package pacta

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestClonePACTAVersion(t *testing.T)               { testClone(t, &PACTAVersion{}) }
func TestCloneUser(t *testing.T)                       { testClone(t, &User{}) }
func TestCloneInitiative(t *testing.T)                 { testClone(t, &Initiative{}) }
func TestCloneInitiativeInvitation(t *testing.T)       { testClone(t, &InitiativeInvitation{}) }
func TestCloneInitiativeUserRelationship(t *testing.T) { testClone(t, &InitiativeUserRelationship{}) }
func TestCloneBlob(t *testing.T)                       { testClone(t, &Blob{}) }
func TestCloneOwner(t *testing.T)                      { testClone(t, &Owner{}) }
func TestCloneIncompleteUpload(t *testing.T)           { testClone(t, &IncompleteUpload{}) }
func TestClonePortfolio(t *testing.T)                  { testClone(t, &Portfolio{}) }
func TestClonePortfolioGroup(t *testing.T)             { testClone(t, &PortfolioGroup{}) }
func TestClonePortfolioGroupMembership(t *testing.T)   { testClone(t, &PortfolioGroupMembership{}) }
func TestClonePortfolioSnapshot(t *testing.T)          { testClone(t, &PortfolioSnapshot{}) }
func TestCloneAnalysis(t *testing.T)                   { testClone(t, &Analysis{}) }
func TestCloneAnalysisArtifact(t *testing.T)           { testClone(t, &AnalysisArtifact{}) }
func TestCloneAuditLog(t *testing.T)                   { testClone(t, &AuditLog{}) }
func TestClonePortfolioInitiativeMembership(t *testing.T) {
	testClone(t, &PortfolioInitiativeMembership{})
}

func testClone[C cloneable[C]](t *testing.T, c C) {
	r := rand.New(rand.NewSource(0))
	t.Helper()
	original := c
	populateStruct(r, original, []string{fmt.Sprintf("%T", original)})
	cloned := original.Clone()

	if diff := cmp.Diff(original, cloned, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}

	checkPointers(t, reflect.ValueOf(original).Elem(), reflect.ValueOf(cloned).Elem())

	var cNil C
	clonedNil := cNil.Clone()
	if clonedNil != cNil {
		t.Errorf("expected nil, got %v", clonedNil)
	}
}

// checkPointers makes sure that our Clone() implementations are doing deep
// copies. Keeping pointers to the same substructs after cloning runs the risk
// of accidentally mutating the original data.
func checkPointers(t *testing.T, v1, v2 reflect.Value) {
	t.Helper()

	switch v1.Kind() {
	case reflect.Ptr:
		if v1.IsNil() || v2.IsNil() {
			t.Errorf("v1 = %v, v2 = %v", v1.Pointer(), v2.Pointer())
			return
		}

		// If the field value is a pointer, check if it's a deep copy
		if v1.Pointer() == v2.Pointer() {
			t.Error("pointer addresses are identical, not a deep copy")
		}
	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			if isFieldUnexported(v1.Type(), i) {
				continue
			}
			field1 := v1.Field(i)
			field2 := v2.Field(i)
			// If the field value is a struct, recursively check its fields
			checkPointers(t, field1, field2)
		}
	case reflect.Slice:
		for j := 0; j < v1.Len(); j++ {
			checkPointers(t, v1.Index(j), v2.Index(j))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String, reflect.Bool:
		// We don't care about primitives, they're handled by our standard cmp.Diff check.
	}
}

func isFieldUnexported(typ reflect.Type, idx int) bool {
	structField := typ.Field(idx)
	return !unicode.IsUpper(rune(structField.Name[0]))
}

func populateStruct(r *rand.Rand, s interface{}, seenTypes []string) {
	value := reflect.ValueOf(s).Elem()
	ts := fmt.Sprintf("%T", s)
	newST := append([]string{ts}, seenTypes...)
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		populateField(r, fieldValue, newST)
	}
}

func populateField(r *rand.Rand, fieldValue reflect.Value, seenTypes []string) {
	if !fieldValue.CanSet() {
		return
	}
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		fieldValue = fieldValue.Elem()
	}

	switch fieldValue.Kind() {
	case reflect.Struct:
		ts := fmt.Sprintf("%T", fieldValue.Addr().Interface())
		// panic("oh no")
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			fieldValue.Set(reflect.ValueOf(getRandomTime(r)))
		}
		seen := false
		for _, st := range seenTypes {
			if ts == st {
				seen = true
			}
		}
		if !seen {
			populateStruct(r, fieldValue.Addr().Interface(), seenTypes)
		}
	case reflect.Map:
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
		}
		populateMap(r, fieldValue, seenTypes)
	case reflect.Ptr:
		populateField(r, fieldValue.Elem(), seenTypes)
	default:
		setRandomValue(r, fieldValue, seenTypes)
	}
}

func setRandomValue(r *rand.Rand, fieldValue reflect.Value, seenTypes []string) {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(fmt.Sprintf("%d", r.Int()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue.SetInt(r.Int63())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue.SetUint(r.Uint64())
	case reflect.Float32, reflect.Float64:
		fieldValue.SetFloat(r.Float64())
	case reflect.Bool:
		fieldValue.SetBool(true)
	case reflect.Slice:
		setRandomSliceValue(r, fieldValue, seenTypes)
	default:
		panic(fmt.Sprintf("unsupported field value: %+v %+v", fieldValue.Kind(), fieldValue))
	}
}

func setRandomSliceValue(r *rand.Rand, fieldValue reflect.Value, seenTypes []string) {
	length := r.Intn(10) + 1
	slice := reflect.MakeSlice(fieldValue.Type(), length, length)

	for i := 0; i < length; i++ {
		populateField(r, slice.Index(i), seenTypes)
	}

	fieldValue.Set(slice)
}

func populateMap(r *rand.Rand, mapValue reflect.Value, seenTypes []string) {
	mapKeyType := mapValue.Type().Key()
	mapValueType := mapValue.Type().Elem()

	length := r.Intn(10) + 1

	for i := 0; i < length; i++ {
		key := reflect.New(mapKeyType).Elem()
		value := reflect.New(mapValueType).Elem()

		populateField(r, key, seenTypes)
		populateField(r, value, seenTypes)

		mapValue.SetMapIndex(key, value)
	}
}

func getRandomTime(r *rand.Rand) time.Time {
	year := r.Intn(20) + 2000
	month := time.Month(r.Intn(12) + 1)
	day := r.Intn(28) + 1
	hour := r.Intn(24)
	minute := r.Intn(60)
	second := r.Intn(60)
	nanosecond := r.Intn(1000000000)

	return time.Date(year, month, day, hour, minute, second, nanosecond, time.UTC)
}
