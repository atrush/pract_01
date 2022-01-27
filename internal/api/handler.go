package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc     service.URLShortener
	baseURL string
}

func NewHandler(svc service.URLShortener, baseURL string) *Handler {

	return &Handler{
		svc:     svc,
		baseURL: baseURL,
	}
}

func (h *Handler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		h.unsupportedMediaTypeError(w, "Разрешены запросы только в формате JSON!")
		return
	}

	jsBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.badRequestError(w, err.Error())
		return
	}
	defer r.Body.Close()

	incoming := ShortenRequest{}
	if err := json.Unmarshal(jsBody, &incoming); err != nil {
		h.badRequestError(w, "неверный формат JSON")
		return
	}
	if err := incoming.Validate(); err != nil {
		h.badRequestError(w, err.Error())
		return
	}

	shortID, err := h.svc.SaveURL(incoming.SrcURL)
	if err != nil {
		h.badRequestError(w, err.Error())
		return
	}

	jsResult, err := json.Marshal(ShortenResponse{
		Result: h.baseURL + "/" + shortID,
	})
	if err != nil {
		h.serverError(w, err.Error())
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsResult)
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	srcURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.serverError(w, err.Error())
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
	w.Write([]byte(h.baseURL + "/" + shortID))
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")
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
func (h *Handler) serverError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusInternalServerError)
}

func (h *Handler) badRequestError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusBadRequest)
}

func (h *Handler) unsupportedMediaTypeError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusUnsupportedMediaType)
}

func (h *Handler) notFoundError(w http.ResponseWriter) {
	http.Error(w, "запрашиваемая страница не найдена", http.StatusNotFound)
}
