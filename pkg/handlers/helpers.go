package handlers

import (
	"math/rand"
	"net/http"

	"github.com/pkg/errors"
)

func GetUserID(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := store.Get(r, "hndaily")
	if err != nil {
		return "", errors.Wrap(err, "failed to get session \"hndaily\" from request")
	}

	userIDVal := session.Values["user_id"]
	userID, ok := userIDVal.(string)
	if !ok {
		userID = generateUserID()
		session.Values["user_id"] = userID
		session.Save(r, w)
	}

	return userID, nil
}

func generateUserID() string {
	b := make([]byte, userIDLength)
	for i := range b {
		b[i] = userIDChars[rand.Intn(len(userIDChars))]
	}
	return string(b)
}
