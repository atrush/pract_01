package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/atrush/pract_01.git/internal/model"
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
				storage.User().AddUser(&model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&model.ShortURL{
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
				storage.User().AddUser(&model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&model.ShortURL{
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

	// defer func() {
	// 	tstSt.DropAll()
	// 	tstSt.Close()
	// }()

	tstSt, err := infile.NewFileStorage("")
	require.NoError(t, err)

	svcSht, err := service.NewShortURLService(tstSt)
	require.NoError(t, err)

	svcUser, err := service.NewUserService(tstSt)
	require.NoError(t, err)

	h, err := NewHandler(svcSht, svcUser, "http://localhost:8080")
	require.NoError(t, err)

	g := &errgroup.Group{}
	for j := 0; j < 20; j++ {
		g.Go(func() error {
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

			if err := json.NewEncoder(buf).Encode(&arrTst); err != nil {
				return err
			}

			request := httptest.NewRequest("POST", "/api/shorten/batch", buf)

			request.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			res.Body.Close()

			arrResp := make([]BatchResponse, 0, 500)
			if err := json.Unmarshal(resBody, &arrResp); err != nil {
				log.Printf("err response: %q", string(resBody))
				return err
			}

			found := 0
			for _, tstEl := range arrTst {
				for _, respEl := range arrResp {
					if tstEl.ID == respEl.ID {
						found++
						continue
					}
				}
			}
			if found != len(arrTst) {
				return fmt.Errorf("в полученных ссылках не найдено %v ссылок", len(arrTst)-found)
			}
			return nil
		})
	}
	err = g.Wait()
	require.NoError(t, err)
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
				storage.URL().SaveURL(&model.ShortURL{
					ID:        uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID:   "1xQ6p+JI",
					URL:       "https://practicum.yandex.ru/",
					IsDeleted: false,
					UserID:    uuid.Nil})
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
				storage.User().AddUser(&model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&model.ShortURL{
					ID:        uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID:   "1xQ6p+JI",
					URL:       "https://practicum.yandex.ru/",
					IsDeleted: false,
					UserID:    uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
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
				storage.User().AddUser(&model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(&model.ShortURL{
					ID:        uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID:   "1xQ6p+JI",
					IsDeleted: false,
					URL:       "https://yandex.ru/",
					UserID:    uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
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

func (tt *HandlerTest) CheckTest(db st.Storage, t *testing.T) {
	svcSht, err := service.NewShortURLService(db)
	require.NoError(t, err)

	svcUser, err := service.NewUserService(db)
	require.NoError(t, err)

	h, err := NewHandler(svcSht, svcUser, "http://localhost:8080")
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
