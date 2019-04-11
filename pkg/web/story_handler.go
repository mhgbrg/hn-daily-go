package web

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/mhgbrg/hndaily/pkg/repo"
)

func GetStory(db *sql.DB, storyRepo repo.StoryRepo) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := mux.Vars(r)["id"]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			return NotFoundError(err)
		}

		story, err := storyRepo.LoadStory(db, storyID)
		if err != nil {
			return NotFoundError(err)
		}

		http.Redirect(w, r, story.URL.String(), http.StatusFound)

		return nil
	}
}

func ReadStory(db *sql.DB, storyRepo repo.StoryRepo, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := mux.Vars(r)["id"]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			return NotFoundError(err)
		}

		story, err := storyRepo.LoadStory(db, storyID)
		if err != nil {
			return NotFoundError(err)
		}

		user, err := sessionStorage.GetOrSetUser(db, w, r)
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

type OKResponse struct{}

func MarkStoryAsRead(db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := mux.Vars(r)["id"]
		storyID, err := strconv.Atoi(storyIDStr)
		if err != nil {
			return NotFoundError(err)
		}

		user, err := sessionStorage.GetOrSetUser(db, w, r)
		if err != nil {
			return InternalServerError(err)
		}

		err = repo.MarkStoryAsRead(db, user.ID, storyID)
		if err != nil {
			return InternalServerError(err)
		}

		contentType, ok := r.Header["Content-Type"]
		if ok && len(contentType) > 0 && contentType[0] == "application/json" {
			res, err := json.Marshal(OKResponse{})
			if err != nil {
				return InternalServerError(err)
			}

			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(res)
			if err != nil {
				return InternalServerError(err)
			}
		} else {
			redirectURLs := r.URL.Query()["redirect-url"]
			var redirectURL string
			if len(redirectURLs) == 1 {
				redirectURL = redirectURLs[0]
			} else {
				redirectURL = "/"
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)
		}

		return nil
	}
}
