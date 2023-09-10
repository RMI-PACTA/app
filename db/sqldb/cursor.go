package sqldb

import (
	"fmt"
	"strconv"

	"github.com/RMI/pacta/db"
)

func offsetToCursor(i int) db.Cursor {
	return db.Cursor(fmt.Sprintf("%d", i))
}

func offsetFromCursor(c db.Cursor) (int, error) {
	if c == "" {
		return 0, nil
	}
	result, err := strconv.Atoi(string(c))
	if err != nil {
		return 0, fmt.Errorf("converting cursor to offset failed for %q: %w", c, err)
	}
	return result, nil
}
