package api

import (
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/atrush/pract_01.git/internal/service"
)

type Handler struct {
	svc service.URLShortener
}

func NewHandler(svc service.URLShortener) *Handler {

	return &Handler{
		svc: svc,
	}
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	srcURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.badRequestError(w, err.Error())
		return
	}
	defer r.Body.Close()

	if string(srcURL) == "" {
		h.badRequestError(w, "нельзя сохранить пустую ссылку")
		return
	}

	shortID, err := h.svc.SaveURL(string(srcURL))
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

	longURL, err := h.svc.GetURL(shortID)
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

func (h *Handler) badRequestError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusBadRequest)
}

func (h *Handler) notFoundError(w http.ResponseWriter) {
	http.Error(w, "запрашиваемая страница не найдена", http.StatusNotFound)
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
