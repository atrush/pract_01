package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/shterrors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	auth    Auth
	svc     service.URLShortener
	baseURL string
}

//  NewHandler init new handler object and return pointer.
func NewHandler(shtSvc service.URLShortener, authSvc service.UserManager, baseURL string) (*Handler, error) {
	return &Handler{
		svc:     shtSvc,
		baseURL: baseURL,
		auth:    NewAuth(authSvc),
	}, nil
}

//  Ping handler check db connection.
//  Return status 200 if db is active.
//  Return status 500 if db not active.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Ping(r.Context()); err != nil {
		h.serverError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

//  DeleteBatch handler async soft remove list urls for user.
//	Return status 202 if list of urls accepted to delete.
func (h *Handler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	var batch BatchDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == uuid.Nil {
		h.serverError(w, "User ID is empty")

		return
	}

	if err := h.svc.DeleteURLList(userID, batch...); err != nil {
		h.serverError(w, err.Error())

		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
}

// SaveBatch handler save list of urls and return list of shorten urls.
// Accept list of pairs [id, url] in json format, model BatchRequest.
// Return status 201 and list of pairs [id, short url] in json format, model BatchResponse, if processed.
func (h *Handler) SaveBatch(w http.ResponseWriter, r *http.Request) {

	//read batch
	var batch []BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//make map id[url] to add
	listToAdd := make(map[string]string, len(batch))
	for _, batchEl := range batch {
		listToAdd[batchEl.ID] = batchEl.URL
	}

	userID := h.getUserIDFromContext(r)

	//save mp to db, values in map updates to shortURL
	if err := h.svc.SaveURLList(listToAdd, userID); err != nil {
		h.serverError(w, err.Error())

		return
	}

	//serialize response
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "   ")
	if err := encoder.Encode(NewBatchListResponseFromMap(listToAdd, h.baseURL)); err != nil {
		h.serverError(w, err.Error())

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buffer.Bytes())
}

// GetUserUrls handler return list of stored urls for user.
// Return status 200 and list of [short url, url] in json format, model ShortenListResponse,  if user urls founded.
// Return status 204 if urls not founded.
func (h *Handler) GetUserUrls(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)

	if userID == uuid.Nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	urlList, err := h.svc.GetUserURLList(r.Context(), userID)
	if err != nil {
		h.serverError(w, err.Error())
		return
	}

	if urlList == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if len(urlList) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsResult, err := json.Marshal(NewShortenListResponseFromCanonical(urlList, h.baseURL))
	if err != nil {
		h.serverError(w, err.Error())
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsResult))
}

// SaveURLJSONHandler save incoming url and return short url.
// Accept url in json format, model ShortenRequest.
// Return status 201 and short url in json format, model ShortenResponse, if saved.
// Return status 409 and stored short url in json format, model ShortenResponse, if url is exist in db.
func (h *Handler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	// read incoming ShortenRequest

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

	// save to db and get shortID
	userID := h.getUserIDFromContext(r)
	shortID, err := h.svc.SaveURL(r.Context(), incoming.SrcURL, userID)

	// handle conflict Add
	isConflict := false
	if err != nil {
		shortID, isConflict = processConflictErr(err)
		if !isConflict {
			h.badRequestError(w, err.Error())

			return
		}
	}

	// marshal response
	jsResult, err := json.Marshal(ShortenResponse{
		Result: h.baseURL + "/" + shortID,
	})
	if err != nil {
		h.serverError(w, err.Error())
		return
	}

	w.Header().Set("content-type", "application/json")
	// write response
	if isConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(jsResult)
}

// SaveURLHandler save incoming url and return short url.
// Accept url in text format.
// Return status 201 and short url in text format if saved.
// Return status 409 and stored short url in text format, if url is exist in db.
func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	// read incoming URL
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
	shortID, err := h.svc.SaveURL(r.Context(), string(srcURL), userID)

	// handle conflict Add
	isConflict := false
	if err != nil {
		shortID, isConflict = processConflictErr(err)
		if !isConflict {
			h.badRequestError(w, err.Error())

			return
		}
	}

	// write response
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	if isConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(h.baseURL + "/" + shortID))
}

// GetURLHandler return redirect for url by incoming shortID param.
// Accept shortID from route params.
// Return status 307 and Location field with stored url in header, if short url founded.
// Return status 410 if short url founded, but mark as deleted.
// Return status 404 if short url not founded.
func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")
	if shortID == "" {
		h.badRequestError(w, "короткая ссылка не может быть пустой")
		return
	}

	storedURL, err := h.svc.GetURL(r.Context(), shortID)
	if err != nil {
		h.badRequestError(w, err.Error())
		return
	}

	if storedURL == (model.ShortURL{}) {
		h.notFoundError(w)
		return
	}

	if storedURL.IsDeleted {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", storedURL.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

// getUserIDFromContext gets user UUID from context
func (h *Handler) getUserIDFromContext(r *http.Request) uuid.UUID {
	ctxID := r.Context().Value(ContextKeyUserID)
	if ctxID == nil {
		return uuid.Nil
	}

	userUUID, err := uuid.Parse(ctxID.(string))
	if err == nil {
		return userUUID
	}

	return uuid.Nil
}

// processConflictErr check if error is conflict adding URL.
// If incoming url exist return shortID from error and true.
func processConflictErr(err error) (string, bool) {
	if err != nil && errors.Is(err, &shterrors.ErrorConflictSaveURL{}) {
		conflictErr, _ := err.(*shterrors.ErrorConflictSaveURL)
		return conflictErr.ExistShortURL, true
	}
	return "", false
}

func (h *Handler) serverError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusInternalServerError)
}

func (h *Handler) badRequestError(w http.ResponseWriter, errText string) {
	http.Error(w, errText, http.StatusBadRequest)
}

func (h *Handler) notFoundError(w http.ResponseWriter) {
	http.Error(w, "запрашиваемая страница не найдена", http.StatusNotFound)
}
