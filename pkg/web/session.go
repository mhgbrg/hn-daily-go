package web

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

type SessionStorage struct {
	Store sessions.Store
	DB    *sql.DB
}

func CreateSessionStorage(db *sql.DB, keys CryptoKeys) SessionStorage {
	store := sessions.NewCookieStore(keys.HashKey, keys.EncryptionKey)
	return SessionStorage{
		DB:    db,
		Store: store,
	}
}

func (sessionStorage *SessionStorage) getSession(r *http.Request) (*sessions.Session, error) {
	session, err := sessionStorage.Store.Get(r, "hndaily")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session \"hndaily\" from request")
	}
	return session, nil
}

func (sessionStorage *SessionStorage) SetUser(w http.ResponseWriter, r *http.Request, user models.User) error {
	session, err := sessionStorage.getSession(r)
	if err != nil {
		return err
	}
	session.Values["user_id"] = user.ExternalID
	err = session.Save(r, w)
	if err != nil {
		return errors.Wrap(err, "failed to save userID to session")
	}
	return nil
}

var ErrUserNotSet = errors.New("user has not been set on session")

func (sessionStorage *SessionStorage) GetUser(r *http.Request) (models.User, error) {
	session, err := sessionStorage.getSession(r)
	if err != nil {
		return models.User{}, err
	}
	sessionUserID, ok := session.Values["user_id"]
	if !ok {
		return models.User{}, ErrUserNotSet
	}
	externalUserID, ok := sessionUserID.(string)
	if !ok {
		return models.User{}, errors.Errorf("failed to cast value %v to string", sessionUserID)
	}
	user, err := repo.LoadUserByExternalID(sessionStorage.DB, externalUserID)
	if err != nil {
		return models.User{}, errors.WithMessagef(err, "failed to load user with externalID=%s", externalUserID)
	}
	return user, nil
}

func (sessionStorage *SessionStorage) AddFlash(w http.ResponseWriter, r *http.Request, flash string) error {
	session, err := sessionStorage.getSession(r)
	if err != nil {
		return err
	}
	session.AddFlash(flash)
	err = session.Save(r, w)
	if err != nil {
		return errors.Wrap(err, "failed to save flash to session")
	}
	return nil
}

func (sessionStorage *SessionStorage) Flashes(w http.ResponseWriter, r *http.Request) ([]string, error) {
	session, err := sessionStorage.getSession(r)
	if err != nil {
		return []string{}, err
	}
	flashes := session.Flashes()
	flashesStr := make([]string, len(flashes))
	for i, flash := range flashes {
		str, ok := flash.(string)
		if !ok {
			return []string{}, errors.Errorf("failed to cast flash message %s to string", flash)
		}
		flashesStr[i] = str
	}
	err = session.Save(r, w)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to remove flashes from session")
	}
	return flashesStr, nil
}

func GetOrSetUser(sessionStorage SessionStorage, w http.ResponseWriter, r *http.Request) (models.User, error) {
	user, err := sessionStorage.GetUser(r)
	if err == ErrUserNotSet {
		user, err = repo.CreateUser(sessionStorage.DB)
		if err != nil {
			return models.User{}, errors.WithMessage(err, "failed to create user")
		}
		err = sessionStorage.SetUser(w, r, user)
		if err != nil {
			return models.User{}, errors.WithMessage(err, "failed to set user on session")
		}
	} else if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to get user from session")
	}
	return user, nil
}
