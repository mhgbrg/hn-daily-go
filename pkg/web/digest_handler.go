package web

import (
	"database/sql"
	"net/http"
	"time"

	templatelib "html/template"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

func GetDigest(templates *Templates, db *sql.DB) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		dateStr := r.URL.Path[len("/digest/"):]
		date, err := models.ParseDate(dateStr)
		if err != nil {
			return NotFoundError(err)
		}

		digest, err := repo.LoadDigest(db, date)
		if err == repo.DigestNotFoundError {
			return NotFoundError(err)
		} else if err != nil {
			return InternalServerError(err)
		}

		return renderPage(templates.Digest, db, w, r, digest)
	}
}

func GetLatestDigest(templates *Templates, db *sql.DB) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		digest, err := repo.LoadLatestDigest(db)
		if err != nil {
			return InternalServerError(err)
		}
		return renderPage(templates.Digest, db, w, r, digest)
	}
}

func renderPage(template *templatelib.Template, db *sql.DB, w http.ResponseWriter, r *http.Request, digest models.Digest) error {
	userID, err := GetOrSetUserID(w, r)
	if err != nil {
		return InternalServerError(err)
	}

	storyIDs := make([]int, len(digest.Stories))
	for i, story := range digest.Stories {
		storyIDs[i] = story.ID
	}

	storyReadMap, err := repo.HasReadStories(db, userID, storyIDs)
	if err != nil {
		return InternalServerError(err)
	}

	viewData := createDigestViewData(digest, storyReadMap)
	err = template.Execute(w, viewData)
	if err != nil {
		return InternalServerError(err)
	}

	return nil
}

type digestViewData struct {
	Year        int
	Month       string
	Day         int
	Weekday     string
	ArchiveURL  string
	GeneratedAt time.Time
	Stories     []digestViewStory
}

type digestViewStory struct {
	Rank        int
	Title       string
	URL         string
	Site        string
	Points      int
	NumComments int
	CommentsURL string
	IsRead      bool
}

func createDigestViewData(digest models.Digest, storyReadMap map[int]bool) digestViewData {
	viewStories := make([]digestViewStory, len(digest.Stories))
	for i, story := range digest.Stories {
		viewStories[i] = digestViewStory{
			Rank:        i + 1,
			Title:       story.Title,
			URL:         StoryURL(story.ID),
			Site:        story.URL.Hostname(),
			Points:      story.Points,
			NumComments: story.NumComments,
			CommentsURL: CommentsURL(story.ExternalID),
			IsRead:      storyReadMap[story.ID],
		}
	}

	return digestViewData{
		Weekday:     digest.Date.ToTime().Weekday().String(),
		Month:       digest.Date.Month.String(),
		Day:         digest.Date.Day,
		Year:        digest.Date.Year,
		ArchiveURL:  ArchiveURL(digest.Date.ToYearMonth()),
		GeneratedAt: digest.GeneratedAt,
		Stories:     viewStories,
	}
}
