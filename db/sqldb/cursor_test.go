package sqldb

import (
	"fmt"
	"testing"
)

func TestOffsetToCursorRoundTrip(t *testing.T) {
	for _, offset := range []int{0, 8, 10000} {
		t.Run(fmt.Sprintf("offset %d", offset), func(t *testing.T) {
			offset := 0
			cursor := offsetToCursor(offset)
			roundTrip, err := offsetFromCursor(cursor)
			if err != nil {
				t.Fatal(err)
			}
			if offset != roundTrip {
				t.Fatalf("offset %d != roundTrip %d", offset, roundTrip)
			}
		})
	}
}
