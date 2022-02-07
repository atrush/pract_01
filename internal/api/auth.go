package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/atrush/pract_01.git/pkg"
)

type (
	Auth struct {
		crypt pkg.Crypt
	}
	contextKey string
)

var (
	ContextKeyUserID = contextKey("user-id")
)

func NewAuth() *Auth {
	return &Auth{
		crypt: *pkg.NewCrypt(),
	}
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := a.authUser(w, r)
		if err != nil {
			http.Error(w, "Ошибка установки ключа пользователя: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, string(userID))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *Auth) authUser(w http.ResponseWriter, r *http.Request) (string, error) {
	if cookie, errCookie := r.Cookie("token"); errCookie == nil {

		token := make([]byte, hex.DecodedLen(len(cookie.Value)))
		_, err := hex.Decode(token, []byte(cookie.Value))
		if err != nil {
			return "", err
		}
		userID, err := a.crypt.DecodeToken(token)
		if err == nil {
			return string(userID), nil
		}
	}

	newUserID, newUserToken, err := a.newUser()
	if err != nil {
		return "", fmt.Errorf("ошибка гнерации нового пользователя: %w", err)
	}

	newCookie := http.Cookie{
		Name:     "token",
		Value:    newUserToken,
		MaxAge:   int((time.Hour * 24 * 30).Seconds()),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &newCookie)

	return newUserID, nil
}

func (a *Auth) newUser() (string, string, error) {
	newUserID, err := a.genUserID()
	if err != nil {
		return "", "", err
	}

	newToken, err := a.crypt.EncodeToken(newUserID)
	if err != nil {
		return "", "", err
	}

	return newUserID, string(toHEX(newToken)), nil
}

func (a *Auth) genUserID() (string, error) {
	rand, err := a.crypt.GenerateRandom(16)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации ID пользователя:%w", err)
	}

	hexToken := make([]byte, hex.EncodedLen(len(rand)))
	hex.Encode(hexToken, rand)

	return string(toHEX(rand)[:16]), nil
}

func toHEX(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return dst
}
