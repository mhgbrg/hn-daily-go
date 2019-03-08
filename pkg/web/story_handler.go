package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/mhgbrg/hndaily/pkg/repo"
)

func ReadStory(db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		storyIDStr := r.URL.Path[len("/story/"):]
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
