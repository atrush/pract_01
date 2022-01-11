package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atrush/pract_01.git/internal/storage/mapstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SaveURLHandler(t *testing.T) {

	tests := []struct {
		name     string
		body     string
		wantCode int
	}{
		{
			name:     "POST empty URL",
			body:     "",
			wantCode: 400,
		},
		{
			name:     "POST URL",
			body:     "https://practicum.yandex.ru/",
			wantCode: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{db: mapstorage.NewStorage()}

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))

			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.SaveURLHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.True(t, res.StatusCode == tt.wantCode, "Ожидался код ответа %d, получен %d", tt.wantCode, w.Code)
		})
	}
}

func TestHandler_GetURLHandler(t *testing.T) {
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name     string
		request  string
		wantCode int
	}{
		{
			name:     "GET not exist URL with short URL",
			request:  "/aaaaaa",
			wantCode: 404,
		},
		{
			name:     "GET not exist URL",
			request:  "/aaaaaa/ddddd",
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{db: mapstorage.NewStorage()}

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(handler.GetURLHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			assert.True(t, res.StatusCode == tt.wantCode, "Ожидался код ответа %d, получен %d", tt.wantCode, w.Code)
		})
	}
}

func Test_testSaveAndGetURL(t *testing.T) {
	longURL := "https://practicum.yandex.ru/"
	longURLHeader := "Location"

	handler := Handler{db: mapstorage.NewStorage()}
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(longURL)))
	w := httptest.NewRecorder()
	h := http.HandlerFunc(handler.SaveURLHandler)
	h.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	shortURL := string(resBody)

	request = httptest.NewRequest(http.MethodGet, shortURL, nil)
	w = httptest.NewRecorder()
	h = http.HandlerFunc(handler.GetURLHandler)
	h.ServeHTTP(w, request)

	res = w.Result()
	defer res.Body.Close()

	assert.True(t, res.StatusCode == 307, "При получении ссылки ожидался код ответа %d, получен %d", 307, w.Code)

	headLocationVal, ok := res.Header[longURLHeader]
	require.True(t, ok, "При получении ссылки не получен заголовок %v", longURLHeader)
	assert.Equal(t, longURL, headLocationVal[0], "Поучена ссылка %v, ожидалась %v", headLocationVal[0], longURL)

}
