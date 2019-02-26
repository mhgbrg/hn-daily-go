package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/mhgbrg/hndaily/pkg/repo"
)

const userIDLength = 6
const userIDChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func ReadStory(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storyIDStr := r.URL.Path[len("/story/"):]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			ReturnError(w, err, 404)
			return
		}

		story, err := repo.LoadStory(db, storyID)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}

		userID, err := GetUserID(w, r)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}

		err = repo.MarkStoryAsRead(db, userID, storyID)
		if err != nil {
			ReturnError(w, err, 500)
			return
		}

		http.Redirect(w, r, story.URL.String(), 302)
	}
}
