package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

type Handler struct {
	db storage.URLStorer
}

func NewHandler(db storage.URLStorer) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	srcURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.badRequestError(w, err.Error())

		return
	}

	if string(srcURL) == "" {
		h.badRequestError(w, "нельзя сохранить пустую ссылку")

		return
	}

	shortID, err := h.genShortURL(string(srcURL), 0, "")
	if err != nil {
		h.badRequestError(w, err.Error())

		return
	}

	_, err = h.db.SaveURL(shortID, string(srcURL))
	if err != nil {
		h.badRequestError(w, err.Error())

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shortID))
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortID := trimFirstRune(r.URL.Path)
	if shortID == "" {
		h.badRequestError(w, "короткая ссылка не может быть пустой")

		return
	}

	longURL, err := h.db.GetURL(shortID)
	if err != nil {
		h.badRequestError(w, err.Error())

		return
	}

	if longURL == "" {
		h.notFoundError(w)

		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) genShortURL(srcURL string, iterationCount int, salt string) (string, error) {
	shortID := service.GenerateShortLink(srcURL, salt)
	if !h.db.IsAvailableID(shortID) {
		iterationCount++
		salt := uuid.New().String()

		var err error
		shortID, err = h.genShortURL(srcURL, iterationCount, salt)
		if err != nil || iterationCount > 10 {

			return "", errors.New("ошибка генерации короткой ссылки")
		}

		return shortID, nil
	}

	return shortID, nil
}

func (h *Handler) badRequestError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusBadRequest)
}

func (h *Handler) notFoundError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)

	return s[i:]
}
