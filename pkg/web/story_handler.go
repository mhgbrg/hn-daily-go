package web

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/mhgbrg/hndaily/pkg/repo"
)

func ReadStory(db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := mux.Vars(r)["id"]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			return NotFoundError(err)
		}

		story, err := repo.LoadStory(db, storyID)
		if err != nil {
			return NotFoundError(err)
		}

		user, err := GetOrSetUser(sessionStorage, w, r)
		if err != nil {
			return InternalServerError(err)
		}

		err = repo.MarkStoryAsRead(db, user.ID, storyID)
		if err != nil {
			return InternalServerError(err)
		}

		http.Redirect(w, r, story.URL.String(), http.StatusFound)

		return nil
	}
}

func MarkStoryAsRead(db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := mux.Vars(r)["id"]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			return NotFoundError(err)
		}

		redirectURLs := r.URL.Query()["redirect-url"]
		var redirectURL string
		if len(redirectURLs) == 1 {
			redirectURL = redirectURLs[0]
		} else {
			redirectURL = "/"
		}

		user, err := GetOrSetUser(sessionStorage, w, r)
		if err != nil {
			return InternalServerError(err)
		}

		err = repo.MarkStoryAsRead(db, user.ID, storyID)
		if err != nil {
			return InternalServerError(err)
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)

		return nil
	}
}

type OKResponse struct{}

func MarkStoryAsReadJSON(db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := mux.Vars(r)["id"]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			return NotFoundError(err)
		}

		user, err := GetOrSetUser(sessionStorage, w, r)
		if err != nil {
			return InternalServerError(err)
		}

		err = repo.MarkStoryAsRead(db, user.ID, storyID)
		if err != nil {
			return InternalServerError(err)
		}

		res, err := json.Marshal(OKResponse{})
		if err != nil {
			return InternalServerError(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(res)
		if err != nil {
			return InternalServerError(err)
		}

		return nil
	}
}
