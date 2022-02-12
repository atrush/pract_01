package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/google/uuid"
)

type (
	Auth struct {
		crypt AuthCrypt
		svc   service.Servicer
	}
	contextKey string
)

var (
	ContextKeyUserID = contextKey("user-id")
)

func NewAuth(svc service.Servicer) *Auth {
	return &Auth{
		crypt: *NewAuthCrypt(),
		svc:   svc,
	}
}

// Return auth middleware
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := a.authUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Read uuid from cookie token. If ok and user exist set ctx, else generate new user and set cookie
func (a *Auth) authUser(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	if cookie, errCookie := r.Cookie("token"); errCookie == nil {
		//decode token
		id, err := a.crypt.DecodeToken(cookie.Value)
		if err != nil {
			return uuid.Nil, fmt.Errorf("ошибка установки ключа пользователя:%w", err)
		}

		//check user
		exist, err := a.svc.User().Exist(id)
		if err != nil {
			return uuid.Nil, fmt.Errorf("ошибка установки ключа пользователя:%w", err)
		}

		if exist {
			return id, nil
		}
	}

	newUserUUID, newUserToken, err := a.newUser()
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка установки ключа пользователя:%w", err)
	}

	newCookie := http.Cookie{
		Name:     "token",
		Value:    newUserToken,
		MaxAge:   int((time.Hour * 24 * 30).Seconds()),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &newCookie)

	return newUserUUID, nil
}

// Add new user to storage, return UUID and token
func (a *Auth) newUser() (uuid.UUID, string, error) {
	newUser, err := a.svc.User().AddUser()
	if err == nil {
		token, err := a.crypt.EncodeUUID(newUser.ID)
		if err == nil {

			return newUser.ID, string(token), nil
		}
	}

	return uuid.Nil, "", fmt.Errorf("ошибка гнерации нового пользователя: %w", err)
}
