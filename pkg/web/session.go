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
	store := sessions.NewCookieStore(
		keys.HashKey,
		keys.EncryptionKey,
		[]byte("CHANGE-THIS-KEY-BEFORE-COMMITTING"),
		nil,
	)
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
	if err == repo.ErrUserNotFound {
		return models.User{}, ErrUserNotSet
	} else if err != nil {
		return models.User{}, errors.WithMessagef(err, "failed to load user with externalID=%s", externalUserID)
	}
	return user, nil
}

type FlashType int

const (
	Success FlashType = iota
	Failure
)

type Flash struct {
	Message string
	Type    FlashType
}

func (flash Flash) Success() bool {
	return flash.Type == Success
}

func (sessionStorage *SessionStorage) AddFlash(w http.ResponseWriter, r *http.Request, flash Flash) error {
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

func (sessionStorage *SessionStorage) Flashes(w http.ResponseWriter, r *http.Request) ([]Flash, error) {
	session, err := sessionStorage.getSession(r)
	if err != nil {
		return []Flash{}, err
	}
	flashes := session.Flashes()
	converted := make([]Flash, len(flashes))
	for i, flash := range flashes {
		f, ok := flash.(*Flash)
		if !ok {
			return []Flash{}, errors.Errorf("failed to cast flash message %v to type Flash", flash)
		}
		converted[i] = *f
	}
	err = session.Save(r, w)
	if err != nil {
		return []Flash{}, errors.Wrap(err, "failed to remove flashes from session")
	}
	return converted, nil
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
