package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage/psql"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	auth    *Auth
	svc     service.Servicer
	psDB    *psql.Storage
	baseURL string
}

// Return new handler
func NewHandler(svc service.Servicer, psDB *psql.Storage, baseURL string) (*Handler, error) {
	return &Handler{
		svc:     svc,
		psDB:    psDB,
		baseURL: baseURL,
		auth:    NewAuth(svc),
	}, nil
}

// Check db connection
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if h.psDB == nil {
		h.serverError(w, "база данных не инициирована")
		return
	}

	if err := h.psDB.Ping(); err != nil {
		h.serverError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Return arr of stored urls for current user
func (h *Handler) GetUserUrls(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)

	if userID == uuid.Nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	urlList, err := h.svc.URL().GetUserURLList(userID)
	if err != nil {
		h.serverError(w, err.Error())
		return
	}

	if urlList == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	responseArr := make([]ShortenListResponse, 0, len(urlList))
	for _, v := range urlList {
		responseArr = append(responseArr, ShortenListResponse{
			ShortURL: h.baseURL + "/" + v.ShortID,
			SrcURL:   v.URL,
		})
	}

	jsResult, err := json.Marshal(responseArr)
	if err != nil {
		h.serverError(w, err.Error())
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsResult))
}

// Save URL with JSON request
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

	userID := h.getUserIDFromContext(r)
	shortID, err := h.svc.URL().SaveURL(incoming.SrcURL, userID)
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

// Save URL with body request
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

	userID := h.getUserIDFromContext(r)
	shortID, err := h.svc.URL().SaveURL(string(srcURL), userID)
	if err != nil {
		h.badRequestError(w, err.Error())
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.baseURL + "/" + shortID))
}

// Return stored URL by short URL
func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")
	if shortID == "" {
		h.badRequestError(w, "короткая ссылка не может быть пустой")
		return
	}

	longURL, err := h.svc.URL().GetURL(shortID)
	if err != nil {
		h.badRequestError(w, err.Error())
		return
	}

	if longURL == "" {
		h.notFoundError(w)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Get user UUID from context
func (h *Handler) getUserIDFromContext(r *http.Request) uuid.UUID {
	ctxID := r.Context().Value(ContextKeyUserID).(string)
	userUUID, err := uuid.Parse(ctxID)

	if err == nil {
		return userUUID
	}

	return uuid.Nil
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
