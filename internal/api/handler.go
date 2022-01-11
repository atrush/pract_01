package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage"
)

type Handler struct {
	db storage.URLStorer
}

func NewHandler(db storage.URLStorer) *Handler {
	return &Handler{
		db: db,
	}
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)

	return s[i:]
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	srcURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.badRequestError(w)
	}

	shortID, err := h.genShortURL(string(srcURL), 0, 0)
	if err != nil {
		h.badRequestError(w)

		return
	}

	_, err = h.db.SaveURL(shortID, string(srcURL))
	if err != nil {
		h.badRequestError(w)

		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/"+shortID)
}

func (h *Handler) genShortURL(srcURL string, saltCount int, iterationCount int) (string, error) {
	shortID := service.GenerateShortLink(srcURL, strconv.Itoa(saltCount))

	if !h.db.IsAvailableID(shortID) {
		saltCount++
		iterationCount++
		shortID, err := h.genShortURL(srcURL, saltCount, iterationCount)
		if err != nil || iterationCount > 10 {

			return "", errors.New("ошибка генерации короткой ссылки")
		}

		return shortID, nil
	}

	return shortID, nil
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	q := trimFirstRune(r.URL.Path)
	longURL, err := h.db.GetURL(q)
	if err != nil {
		h.badRequestError(w)

		return
	}
	if longURL == "" {
		h.notFoundError(w)

		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func (h *Handler) BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Разрешены только GET и POST запросы", http.StatusBadRequest)
}

func (h *Handler) badRequestError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (h *Handler) notFoundError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}