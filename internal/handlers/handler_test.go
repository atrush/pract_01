package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_RequestHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code         int
		wantBody     bool
		responseBody string
	}
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		method  string
		name    string
		request string
		body    string
		want    want
	}{
		// определяем все тесты
		{
			name:    "POST empty URL",
			method:  http.MethodPost,
			request: "/",
			body:    "",
			want: want{
				code:     400,
				wantBody: false,
			},
		},
		{
			name:    "POST URL",
			method:  http.MethodPost,
			request: "/",
			body:    "https://practicum.yandex.ru/",
			want: want{
				code:     201,
				wantBody: false,
			},
		},
		{
			name:    "PUT URL",
			method:  http.MethodPut,
			request: "/",
			body:    "https://practicum.yandex.ru/",
			want: want{
				code:     400,
				wantBody: false,
			},
		},
		{
			name:    "GET not exist URL",
			method:  http.MethodGet,
			request: "/aaaaaa",
			want: want{
				code:     404,
				wantBody: false,
			},
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{DB: storage.NewStorage()}

			request := httptest.NewRequest(tt.method, tt.request, bytes.NewBuffer([]byte(tt.body)))

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(handler.RequestHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.True(t, res.StatusCode == tt.want.code, "Ожидался код ответа %d, получен %d", tt.want.code, w.Code)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.False(t, tt.want.wantBody && (tt.want.responseBody != string(resBody)),
				"Ожидалось тело ответа %d, получено %d", tt.want.wantBody, w.Body.String())
		})
	}

	testSaveAndGetUrl(t)
}

func testSaveAndGetUrl(t *testing.T) {
	longURL := "https://practicum.yandex.ru/"
	longURLHeader := "Location"

	handler := Handler{DB: storage.NewStorage()}
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(longURL)))
	w := httptest.NewRecorder()
	h := http.HandlerFunc(handler.RequestHandler)
	h.ServeHTTP(w, request)

	res := w.Result()

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	shortUrl := string(resBody)

	request = httptest.NewRequest(http.MethodGet, shortUrl, nil)
	w = httptest.NewRecorder()
	h = http.HandlerFunc(handler.RequestHandler)
	h.ServeHTTP(w, request)

	res = w.Result()
	assert.True(t, res.StatusCode == 307, "При получении ссылки ожидался код ответа %d, получен %d", 307, w.Code)

	headLocationVal, ok := res.Header[longURLHeader]
	require.True(t, ok, "При получении ссылки не получен заголовок %v", longURLHeader)
	assert.Equal(t, longURL, headLocationVal[0], "Поучена ссылка %v, ожидалась %v", headLocationVal[0], longURL)

}
