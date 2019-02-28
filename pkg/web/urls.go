package web

import (
	"fmt"

	"github.com/mhgbrg/hndaily/pkg/models"
)

func DigestURL(date models.Date) string {
	return fmt.Sprintf("/digest/%s", date)
}

func StoryURL(id int) string {
	return fmt.Sprintf("/story/%d", id)
}

func CommentsURL(externalID int) string {
	return fmt.Sprintf("https://news.ycombinator.com/item?id=%d", externalID)
}

func ArchiveURL(yearMonth models.YearMonth) string {
	return fmt.Sprintf("/archive/%s", yearMonth)
}
