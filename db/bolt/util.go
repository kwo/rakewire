package bolt

import (
	"fmt"
	"rakewire.com/db"
	"time"
)

func fetchKey(f *db.Feed) string {
	return fmt.Sprintf("%s!%s", formatFetchTime(*f.GetNextFetchTime()), f.ID)
}

func formatFetchTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

func formatMaxTime(t time.Time) string {
	return formatFetchTime(t) + "#"
}
