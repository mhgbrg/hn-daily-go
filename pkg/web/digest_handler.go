package web

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	urllib "net/url"
	"time"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

func GetDigest(templates *Templates, db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		dateStr := mux.Vars(r)["date"]
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

		return renderPage(templates, db, sessionStorage, w, r, digest)
	}
}

func GetLatestDigest(templates *Templates, db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		digest, err := repo.LoadLatestDigest(db)
		if err != nil {
			return InternalServerError(err)
		}
		return renderPage(templates, db, sessionStorage, w, r, digest)
	}
}

func SetDeviceID(db *sql.DB, sessionStorage SessionStorage) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		newExternalUserID := r.FormValue("deviceID")

		user, err := repo.LoadUserByExternalID(db, newExternalUserID)
		if err == repo.ErrUserNotFound {
			err = sessionStorage.AddFlash(w, r, Flash{"Invalid device ID", Failure})
			if err != nil {
				return InternalServerError(err)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		} else if err != nil {
			return InternalServerError(err)
		}

		err = sessionStorage.SetUser(w, r, user)
		if err != nil {
			return InternalServerError(err)
		}

		err = sessionStorage.AddFlash(w, r, Flash{"Device ID updated successfully!", Success})
		if err != nil {
			return InternalServerError(err)
		}

		http.Redirect(w, r, "/", http.StatusFound)

		return nil
	}
}

func renderPage(
	templates *Templates,
	db *sql.DB,
	sessionStorage SessionStorage,
	w http.ResponseWriter,
	r *http.Request,
	digest models.Digest,
) error {
	user, err := GetOrSetUser(sessionStorage, w, r)
	if err != nil {
		return InternalServerError(err)
	}

	storyIDs := make([]int, len(digest.Stories))
	for i, story := range digest.Stories {
		storyIDs[i] = story.ID
	}

	storyReadMap, err := repo.HasReadStories(db, user.ID, storyIDs)
	if err != nil {
		return InternalServerError(err)
	}

	flashes, err := sessionStorage.Flashes(w, r)
	if err != nil {
		return InternalServerError(err)
	}

	viewData := createDigestViewData(r.URL, digest, storyReadMap, user, flashes)

	err = templates.Digest.Execute(w, viewData)
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
	UserID      string
	Flashes     []Flash
}

type digestViewStory struct {
	Rank          int
	Title         string
	URL           string
	Site          string
	Points        int
	NumComments   int
	CommentsURL   string
	IsRead        bool
	MarkAsReadURL string
}

func createDigestViewData(
	url *urllib.URL,
	digest models.Digest,
	storyReadMap map[int]bool,
	user models.User,
	flashes []Flash,
) digestViewData {
	viewStories := make([]digestViewStory, len(digest.Stories))
	for i, story := range digest.Stories {
		viewStories[i] = digestViewStory{
			Rank:          i + 1,
			Title:         story.Title,
			URL:           StoryURL(story.ID),
			Site:          story.URL.Hostname(),
			Points:        story.Points,
			NumComments:   story.NumComments,
			CommentsURL:   CommentsURL(story.ExternalID),
			IsRead:        storyReadMap[story.ID],
			MarkAsReadURL: MarkAsReadURL(story.ID, url),
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
		UserID:      user.ExternalID,
		Flashes:     flashes,
	}
}
