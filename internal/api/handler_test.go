package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atrush/pract_01.git/internal/service"
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/atrush/pract_01.git/pkg"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HandlerTest struct {
	name        string
	method      string
	url         string
	body        string
	contentType string

	outBodyExpected     string
	outContTypeExpected string
	outCodeExpected     int

	initFixtures func(storage st.Storage)
}

func TestHandler_SaveConflict(t *testing.T) {
	cfg, err := pkg.NewConfig()
	require.NoError(t, err)
	tests := []HandlerTest{
		{
			name:                "POST conflict /api/shorten",
			method:              http.MethodPost,
			url:                 "/api/shorten",
			body:                "{\"url\": \"https://practicum.yandex.ru\"}",
			contentType:         "application/json",
			outCodeExpected:     409,
			outBodyExpected:     fmt.Sprintf("{\"result\":\"%v/1xQ6p+JI\"}", cfg.BaseURL),
			outContTypeExpected: "application/json",
			initFixtures: func(storage st.Storage) {
				storage.User().AddUser(&st.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&st.ShortURL{
					ID:      uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID: "1xQ6p+JI",
					URL:     "https://practicum.yandex.ru",
					UserID:  uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
			},
		},
		{
			name:                "POST conflict /",
			method:              http.MethodPost,
			url:                 "/",
			body:                "https://practicum.yandex.ru/",
			outCodeExpected:     409,
			outContTypeExpected: "text/plain; charset=utf-8",
			outBodyExpected:     fmt.Sprintf("%v/1xQ6p+JI", cfg.BaseURL),
			initFixtures: func(storage st.Storage) {
				storage.User().AddUser(&st.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&st.ShortURL{
					ID:      uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID: "1xQ6p+JI",
					URL:     "https://practicum.yandex.ru/",
					UserID:  uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//tstSt, err := psql.NewTestStorage("postgres://postgres:hjvfirb@localhost:5432/tst_00?sslmode=disable")
			tstSt, err := infile.NewFileStorage("")
			require.NoError(t, err)
			require.NoError(t, err)

			if tt.initFixtures != nil {
				tt.initFixtures(tstSt)
			}
			tt.CheckTest(tstSt, t)
			// tstSt.DropAll()
			// tstSt.Close()
		})
	}
}
func TestHandler_ShortenBatch(t *testing.T) {

	// tstSt, err := psql.NewTestStorage("postgres://postgres:hjvfirb@localhost:5432/tst_00?sslmode=disable")
	// require.NoError(t, err)

	// defer tstSt.Close()
	// defer tstSt.DropAll()

	tstSt, err := infile.NewFileStorage("")
	require.NoError(t, err)

	svc, err := service.NewService(tstSt)
	require.NoError(t, err)

	h, err := NewHandler(svc, "http://localhost:8080")
	require.NoError(t, err)

	r := NewRouter(h)

	arrTst := make([]BatchRequest, 0, 500)
	for i := 0; i < 500; i++ {
		id := uuid.New().String()
		arrTst = append(arrTst, BatchRequest{
			ID:  uuid.New().String(),
			URL: "http://localhost:8080/" + id,
		})
	}

	buf := new(bytes.Buffer)

	assert.NoError(t, json.NewEncoder(buf).Encode(&arrTst))

	request := httptest.NewRequest("POST", "/api/shorten/batch", buf)

	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	res := w.Result()
	arrResp := make([]BatchResponse, 0, 500)
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&arrResp))
	res.Body.Close()

	found := 0
	for _, tstEl := range arrTst {
		for _, respEl := range arrResp {
			if tstEl.ID == respEl.ID {
				found++
				continue
			}
		}
	}
	require.Equal(t, found, len(arrTst), "в полученных ссылках не найдено %v ссылок", len(arrTst)-found)
}

func TestHandler_SaveURLHandler(t *testing.T) {
	tests := []HandlerTest{
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
			initFixtures: func(storage st.Storage) {
				storage.URL().SaveURL(&st.ShortURL{
					ID:      uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID: "1xQ6p+JI",
					URL:     "https://practicum.yandex.ru/",
					UserID:  uuid.Nil})
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
			url:             "/1xQ6p+JI",
			outCodeExpected: 307,
			initFixtures: func(storage st.Storage) {
				storage.User().AddUser(&st.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&st.ShortURL{
					ID:      uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID: "1xQ6p+JI",
					URL:     "https://practicum.yandex.ru/",
					UserID:  uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
			},
		},
		{
			name:            "POST JSON URL not valid content-type",
			method:          http.MethodPost,
			url:             "/api/shorten",
			body:            "{\"url\": \"https://practicum.yandex.ru/\"}",
			contentType:     "text/plain; charset=utf-8",
			outCodeExpected: 415,
		},
		{
			name:            "POST JSON empty URL",
			method:          http.MethodPost,
			url:             "/api/shorten",
			body:            "{\"url\": \"\"}",
			contentType:     "application/json",
			outCodeExpected: 400,
		},
		{
			name:            "POST JSON empty body",
			method:          http.MethodPost,
			url:             "/api/shorten",
			body:            "",
			contentType:     "application/json",
			outCodeExpected: 400,
		},
		{
			name:            "POST JSON not valid JSON",
			method:          http.MethodPost,
			url:             "/api/shorten",
			body:            "{\"url: \"\"}",
			contentType:     "application/json",
			outCodeExpected: 400,
		},
		{
			name:                "POST JSON URL empty storage",
			method:              http.MethodPost,
			url:                 "/api/shorten",
			body:                "{\"url\": \"https://practicum.yandex.ru/\"}",
			contentType:         "application/json",
			outCodeExpected:     201,
			outContTypeExpected: "application/json",
		},
		{
			name:                "POST JSON URL not empty storage",
			method:              http.MethodPost,
			url:                 "/api/shorten",
			body:                "{\"url\": \"https://practicum.yandex.ru/\"}",
			contentType:         "application/json",
			outCodeExpected:     201,
			outContTypeExpected: "application/json",
			initFixtures: func(storage st.Storage) {
				storage.User().AddUser(&st.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&st.ShortURL{
					ID:      uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID: "1xQ6p+JI",
					URL:     "https://yandex.ru/",
					UserID:  uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := infile.NewFileStorage("")
			require.NoError(t, err)

			if tt.initFixtures != nil {
				tt.initFixtures(db)
			}
			tt.CheckTest(db, t)
		})
	}
}
func Test_testSaveAndGetURL(t *testing.T) {
	cfg, err := pkg.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	longURL := "https://practicum.yandex.ru/"
	longURLHeader := "Location"

	db, err := infile.NewFileStorage("")
	require.NoError(t, err)

	svc, err := service.NewService(db)
	require.NoError(t, err)

	h, err := NewHandler(svc, "http://localhost:8080")
	require.NoError(t, err)

	r := NewRouter(h)
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(longURL)))
	request.RemoteAddr = "localhost" + cfg.ServerPort
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	res := w.Result()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	shortURL := string(resBody)

	request = httptest.NewRequest(http.MethodGet, shortURL, nil)
	request.RemoteAddr = "localhost" + cfg.ServerPort
	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
	res = w.Result()
	defer res.Body.Close()
	assert.True(t, res.StatusCode == 307, "При получении ссылки ожидался код ответа %d, получен %d", 307, w.Code)

	headLocationVal, ok := res.Header[longURLHeader]
	require.True(t, ok, "При получении ссылки не получен заголовок %v", longURLHeader)
	assert.Equal(t, longURL, headLocationVal[0], "Поучена ссылка %v, ожидалась %v", headLocationVal[0], longURL)
}

func (tt *HandlerTest) CheckTest(db st.Storage, t *testing.T) {
	svc, err := service.NewService(db)
	require.NoError(t, err)

	h, err := NewHandler(svc, "http://localhost:8080")
	require.NoError(t, err)

	r := NewRouter(h)
	request := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(tt.body)))

	if tt.contentType != "" {
		request.Header.Set("Content-Type", tt.contentType)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	res := w.Result()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	strBody := string(resBody)

	fmt.Printf("%v: res - %v\n", tt.name, strBody)

	if tt.outContTypeExpected != "" {
		assert.True(t, res.Header.Get("Content-Type") == tt.outContTypeExpected,
			"Ожидался content-type ответа %v, получен %v", tt.outContTypeExpected, res.Header.Get("Content-Type"))
	}

	if tt.outCodeExpected != 0 {
		assert.True(t, res.StatusCode == tt.outCodeExpected, "Ожидался код ответа %d, получен %d", tt.outCodeExpected, w.Code)
	}

	if tt.outBodyExpected != "" {
		assert.Equal(t, strBody, tt.outBodyExpected, "Ожидался ответа %v, получен %v", tt.outBodyExpected, strBody)
	}

}
