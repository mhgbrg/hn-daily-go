package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	templatelib "html/template"
	"log"
	"net/http"
	"time"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

func GetDigest(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	template, err := templatelib.ParseFiles("templates/digest.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Path[len("/digest/"):]
		date, err := models.ParseDate(dateStr)
		if err != nil {
			ReturnError(w, err, 404)
			return
		}

		digest, err := repo.LoadDigest(db, date)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}

		userID, err := GetUserID(w, r)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}

		storyIDs := make([]int, len(digest.Stories))
		for i, story := range digest.Stories {
			storyIDs[i] = story.ID
		}

		storyReadMap, err := repo.HasReadStories(db, userID, storyIDs)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}

		viewData := createViewData(digest, storyReadMap)
		var responseBody bytes.Buffer
		err = template.Execute(&responseBody, viewData)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}
		fmt.Fprint(w, responseBody.String())
	}
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
