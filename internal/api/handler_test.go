package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SaveURLHandler(t *testing.T) {

	tests := []struct {
		name        string
		method      string
		url         string
		body        string
		contentType string

		outBodyExpected     string
		outContTypeExpected string
		outCodeExpected     int

		initFixtures func(storage *inmemory.MapStorage)
	}{
		{
			name:            "POST empty URL",
			method:          http.MethodPost,
			url:             "/",
			body:            "",
			outCodeExpected: 400,
		},
		{
			name:            "POST URL",
			method:          http.MethodPost,
			url:             "/",
			body:            "https://practicum.yandex.ru/",
			outCodeExpected: 201,
		},
		{
			name:            "POST exist URL",
			method:          http.MethodPost,
			url:             "/",
			body:            "https://practicum.yandex.ru/",
			outCodeExpected: 201,
			initFixtures: func(storage *inmemory.MapStorage) {
				storage.SaveURL("wdwfff", "https://practicum.yandex.ru/")
			},
		},
		{
			name:            "GET not exist URL",
			method:          http.MethodGet,
			url:             "/aaaaaa",
			outCodeExpected: 404,
		},
		{
			name:            "GET exist URL",
			method:          http.MethodGet,
			url:             "/aaaaaa",
			outCodeExpected: 307,
			initFixtures: func(storage *inmemory.MapStorage) {
				storage.SaveURL("aaaaaa", "https://practicum.yandex.ru/")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmemory.NewStorage()
			if tt.initFixtures != nil {
				tt.initFixtures(db)
			}

			svc, _ := service.NewShortURLService(db)
			h := &Handler{svc: svc}
			r := NewRouter(*h)

			request := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.True(t, tt.outCodeExpected == 0 || res.StatusCode == tt.outCodeExpected, "Ожидался код ответа %d, получен %d", tt.outCodeExpected, w.Code)
		})
	}
}

func Test_testSaveAndGetURL(t *testing.T) {
	longURL := "https://practicum.yandex.ru/"
	longURLHeader := "Location"

	db := inmemory.NewStorage()
	svc, _ := service.NewShortURLService(db)
	handler := Handler{svc: svc}
	r := NewRouter(handler)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(longURL)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)

	}

	shortURL := string(resBody)

	request = httptest.NewRequest(http.MethodGet, shortURL, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	res = w.Result()
	defer res.Body.Close()

	assert.True(t, res.StatusCode == 307, "При получении ссылки ожидался код ответа %d, получен %d", 307, w.Code)

	headLocationVal, ok := res.Header[longURLHeader]
	require.True(t, ok, "При получении ссылки не получен заголовок %v", longURLHeader)
	assert.Equal(t, longURL, headLocationVal[0], "Поучена ссылка %v, ожидалась %v", headLocationVal[0], longURL)

}
