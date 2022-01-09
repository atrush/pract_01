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

func (h *Handler) RequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := trimFirstRune(r.URL.Path)
		fmt.Println(q)
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

		return
	}

	if r.Method == http.MethodPost {
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
		return
	}
	http.Error(w, "Only GET and POST requests are allowed!", http.StatusBadRequest)
}

func (h *Handler) badRequestError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (h *Handler) notFoundError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
