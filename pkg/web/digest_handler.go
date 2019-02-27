package web

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

func GetDigest(db *sql.DB) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (fmt.Stringer, error) {
		dateStr := r.URL.Path[len("/digest/"):]
		date, err := models.ParseDate(dateStr)
		if err != nil {
			return nil, NotFoundError(err)
		}

		digest, err := repo.LoadDigest(db, date)
		if err != nil {
			return nil, InternalServerError(err)
		}

		return renderPage(db, w, r, digest)
	}
}

func GetLatestDigest(db *sql.DB) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (fmt.Stringer, error) {
		digest, err := repo.LoadLatestDigest(db)
		if err != nil {
			return nil, InternalServerError(err)
		}
		return renderPage(db, w, r, digest)
	}
}

func renderPage(db *sql.DB, w http.ResponseWriter, r *http.Request, digest models.Digest) (fmt.Stringer, error) {
	userID, err := GetOrSetUserID(w, r)
	if err != nil {
		return nil, InternalServerError(err)
	}

	storyIDs := make([]int, len(digest.Stories))
	for i, story := range digest.Stories {
		storyIDs[i] = story.ID
	}

	storyReadMap, err := repo.HasReadStories(db, userID, storyIDs)
	if err != nil {
		return nil, InternalServerError(err)
	}

	template, err := GetTemplate("digest")
	if err != nil {
		return nil, InternalServerError(err)
	}

	viewData := createViewData(digest, storyReadMap)
	var responseBody bytes.Buffer
	err = template.Execute(&responseBody, viewData)
	if err != nil {
		return nil, InternalServerError(err)
	}

	return &responseBody, nil
}

type ViewData struct {
	Year        int
	Month       string
	Day         int
	Weekday     string
	GeneratedAt time.Time
	Stories     []ViewStory
}

type ViewStory struct {
	Rank        int
	Title       string
	URL         string
	Site        string
	Points      int
	NumComments int
	CommentsURL string
	IsRead      bool
}

func createViewData(digest models.Digest, storyReadMap map[int]bool) ViewData {
	viewData := ViewData{
		Weekday:     digest.Date.ToTime().Weekday().String(),
		Month:       digest.Date.Month.String(),
		Day:         digest.Date.Day,
		Year:        digest.Date.Year,
		GeneratedAt: digest.GeneratedAt,
		Stories:     make([]ViewStory, len(digest.Stories)),
	}

	for i, story := range digest.Stories {
		viewStory := ViewStory{
			Rank:        i + 1,
			Title:       story.Title,
			URL:         fmt.Sprintf("/story/%d", story.ID),
			Site:        story.URL.Hostname(),
			Points:      story.Points,
			NumComments: story.NumComments,
			CommentsURL: fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ExternalID),
			IsRead:      storyReadMap[story.ID],
		}
		viewData.Stories[i] = viewStory
	}

	return viewData
}
