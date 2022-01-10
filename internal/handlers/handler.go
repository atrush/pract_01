package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/atrush/pract_01.git/internal/storage"
)

type Handler struct {
	DB storage.URLStorer
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	longURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.badRequestError(w)
	}
	shortURL, err := h.DB.SaveURL(string(longURL))
	if err != nil {
		h.badRequestError(w)
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/"+shortURL)
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	q := trimFirstRune(r.URL.Path)
	longURL, err := h.DB.GetURL(q)
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
