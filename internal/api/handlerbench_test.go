package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test save batch URL handler
func BenchmarkBatchSaveURL(b *testing.B) {
	r, w, cookie, err := initTestHandler(b)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer() // останавливаем таймер

		// array to save
		arrTst := make([]BatchRequest, 0, 15)
		for i := 0; i < cap(arrTst); i++ {
			id := uuid.New().String()
			arrTst = append(arrTst, BatchRequest{
				ID:  uuid.New().String(),
				URL: "http://localhost:8080/" + id,
			})
		}

		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(&arrTst); err != nil {
			b.Fatal(err)
		}

		b.StartTimer() // возобновляем таймер

		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", buf)
		request.AddCookie(cookie)
		request.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, request)
	}
}

// Test save URL handler with JSON body
func BenchmarkBJSONSaveURL(b *testing.B) {
	r, w, cookie, err := initTestHandler(b)
	if err != nil {
		b.ResetTimer()
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("{\"url\": \"https://practicum.yandex.ru/%v\"}", uuid.New().String())

		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(url)))
		request.AddCookie(cookie)
		request.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, request)
	}
}

// Test save URL handler with text body
func BenchmarkTextSaveURL(b *testing.B) {
	r, w, cookie, err := initTestHandler(b)
	if err != nil {
		b.ResetTimer()
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("https://practicum.yandex.ru/%v", uuid.New().String())

		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(url)))
		request.AddCookie(cookie)
		request.Header.Set("Content-Type", "text/plain; charset=utf-8")

		r.ServeHTTP(w, request)
	}
}

// Get auth cookie with ping request
func getAuthCookie(r *chi.Mux, w *httptest.ResponseRecorder) (*http.Cookie, error) {
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")

	r.ServeHTTP(w, request)
	res := w.Result()
	defer res.Body.Close()

	cookies := res.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("not received auth cookie ")
	}
	return cookies[0], nil
}

// init handler from db
func initTestHandler(b *testing.B) (*chi.Mux, *httptest.ResponseRecorder, *http.Cookie, error) {
	b.StopTimer() // останавливаем таймер

	// init db
	db, err := infile.NewFileStorage("")

	if err != nil {
		return nil, nil, nil, err
	}

	// init url service
	svcSht, err := service.NewShortURLService(db)
	if err != nil {
		return nil, nil, nil, err
	}

	// init user service
	svcUser, err := service.NewUserService(db)
	if err != nil {
		return nil, nil, nil, err
	}

	// init handler
	h, err := NewHandler(svcSht, svcUser, "http://localhost:8080")
	if err != nil {
		return nil, nil, nil, err
	}

	r := NewRouter(h, false)
	w := httptest.NewRecorder()

	cookie, err := getAuthCookie(r, w)
	if err != nil {
		return nil, nil, nil, err
	}
	b.StartTimer() // возобновляем таймер
	return r, w, cookie, nil
}
