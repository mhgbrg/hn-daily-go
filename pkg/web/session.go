package web

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/pkg/errors"
)

var store = sessions.NewCookieStore([]byte("CHANGE-THIS-KEY-BEFORE-COMMITTING"))

func GetOrSetUserID(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := store.Get(r, "hndaily")
	if err != nil {
		return "", errors.Wrap(err, "failed to get session \"hndaily\" from request")
	}

	userIDVal := session.Values["user_id"]
	userID, ok := userIDVal.(string)
	if !ok {
		userID = models.GenerateUserID()
		session.Values["user_id"] = userID
		session.Save(r, w)
	}

	return userID, nil
}
