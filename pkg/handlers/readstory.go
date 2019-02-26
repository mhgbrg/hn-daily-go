package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/mhgbrg/hndaily/pkg"
)

var store = sessions.NewCookieStore([]byte("CHANGE-THIS-KEY-BEFORE-COMMITTING"))

const userIDLength = 6
const userIDChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func ReadStory(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storyIDStr := r.URL.Path[len("/story/"):]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}

		story, err := pkg.GetStory(db, storyID)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}

		userID, err := GetUserID(w, r)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}

		err = pkg.MarkStoryAsRead(db, userID, storyID)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, story.URL.String(), 302)
	}
}
