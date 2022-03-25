package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/service"
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/atrush/pract_01.git/pkg"
	"github.com/go-chi/chi/v5"
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
				storage.User().AddUser(context.Background(), model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(context.Background(), model.ShortURL{
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
				storage.User().AddUser(context.Background(), model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(context.Background(), model.ShortURL{
					ID:      uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID: "1xQ6p+JI",
					URL:     "https://practicum.yandex.ru/",
					UserID:  uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tstSt, err := infile.NewFileStorage("")
			require.NoError(t, err)

			if tt.initFixtures != nil {
				tt.initFixtures(tstSt)
			}
			tt.CheckTest(tstSt, t)
		})
	}
}

type GoSaveBatch struct {
	count  int
	cookie *http.Cookie
	saved  []BatchResponse
	r      *chi.Mux
}

func (g *GoSaveBatch) SaveBatch() error {

	arrTst := make([]BatchRequest, 0, g.count)
	for i := 0; i < g.count; i++ {
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
	g.r.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	cookies := res.Cookies()

	//  save cookie
	g.cookie = cookies[0]

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	g.saved = make([]BatchResponse, 0, g.count)
	if err := json.Unmarshal(resBody, &g.saved); err != nil {
		log.Printf("err response: %q", string(resBody))
		return err
	}

	found := 0
	for _, tstEl := range arrTst {
		for _, respEl := range g.saved {
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
}

// send delete batch
func (g *GoSaveBatch) DeleteBatch() error {
	arrDel := make([]string, 0, len(g.saved))
	for _, v := range g.saved {
		arrDel = append(arrDel, strings.Replace(v.ShortURL, "http://localhost:8080/", "", 1))
	}

	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(&arrDel); err != nil {
		return err
	}

	request := httptest.NewRequest("DELETE", "/api/user/urls", buf)
	request.AddCookie(g.cookie)

	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	g.r.ServeHTTP(w, request)

	res := w.Result()
	res.Body.Close()

	if res.StatusCode != 202 {
		return fmt.Errorf("проверка ответа хендлера удаления, код ответа %v вместо 202", res)
	}

	return nil
}

// check deleted Urls
func (g *GoSaveBatch) CheckDeleted() error {

	for _, v := range g.saved {
		buf := new(bytes.Buffer)
		request := httptest.NewRequest("GET", strings.Replace(v.ShortURL, "http://localhost:8080", "", 1), buf)
		request.AddCookie(g.cookie)

		w := httptest.NewRecorder()
		g.r.ServeHTTP(w, request)

		res := w.Result()
		res.Body.Close()
		if res.StatusCode != 410 {
			return fmt.Errorf("проверка удаленной ссылки, код ответа %v вместо 410", res.StatusCode)
		}
	}
	return nil
}

func TestHandler_BatchDelete(t *testing.T) {
	tstSt, err := infile.NewFileStorage("")
	require.NoError(t, err)

	h := initHandler(t, tstSt)

	g := &errgroup.Group{}
	workersCount := 1

	for j := 0; j < workersCount; j++ {
		g.Go(func() error {
			saveBatchItem := GoSaveBatch{r: NewRouter(h, false), count: 11}
			if err := saveBatchItem.SaveBatch(); err != nil {
				return err
			}

			if err := saveBatchItem.DeleteBatch(); err != nil {
				return err
			}

			time.Sleep(15 * time.Second)

			if err := saveBatchItem.CheckDeleted(); err != nil {
				return err
			}

			return nil
		})
	}
	err = g.Wait()
	require.NoError(t, err)

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

	h := initHandler(t, tstSt)

	g := &errgroup.Group{}
	workersCount := 30

	for j := 0; j < workersCount; j++ {
		g.Go(func() error {
			saveBatchItem := GoSaveBatch{r: NewRouter(h, false), count: 500}
			if err := saveBatchItem.SaveBatch(); err != nil {
				return err
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
			name:            "DELETE List ShortID",
			method:          http.MethodDelete,
			url:             "/api/user/urls",
			body:            "[\"url0\", \"url1\", \"url2\", \"url3\", \"url4\"]",
			contentType:     "application/json",
			outCodeExpected: 202,
			initFixtures: func(storage st.Storage) {
				storage.User().AddUser(context.Background(), model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				for i := 0; i < 5; i++ {
					storage.URL().SaveURL(context.Background(), model.ShortURL{
						ID:        uuid.New(),
						ShortID:   fmt.Sprintf("url%v", i),
						URL:       fmt.Sprintf("https://practicum.yandex.ru/%v", i),
						IsDeleted: false,
						UserID:    uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
				}
			},
		},

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
				storage.URL().SaveURL(context.Background(), model.ShortURL{
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
				storage.User().AddUser(context.Background(), model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(context.Background(), model.ShortURL{
					ID:        uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID:   "1xQ6p+JI",
					URL:       "https://practicum.yandex.ru/",
					IsDeleted: false,
					UserID:    uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1")})
			},
		},
		{
			name:            "GET exist Deleted URL",
			method:          http.MethodGet,
			url:             "/1xQ6p+JI",
			outCodeExpected: 410,
			initFixtures: func(storage st.Storage) {
				storage.User().AddUser(context.Background(), model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(context.Background(), model.ShortURL{
					ID:        uuid.MustParse("49dad1e7-983a-4101-a991-aa0e9523a3b1"),
					ShortID:   "1xQ6p+JI",
					URL:       "https://practicum.yandex.ru/",
					IsDeleted: true,
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
				storage.User().AddUser(context.Background(), model.User{
					ID: uuid.MustParse("34e693a6-78e5-4a2f-a6bb-2fad5da50de1"),
				})
				storage.URL().SaveURL(context.Background(), model.ShortURL{
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
	h := initHandler(t, db)

	r := NewRouter(h, false)
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

// init handler from db
func initHandler(t *testing.T, tstSt st.Storage) *Handler {
	svcSht, err := service.NewShortURLService(tstSt)
	require.NoError(t, err)

	svcUser, err := service.NewUserService(tstSt)
	require.NoError(t, err)

	h, err := NewHandler(svcSht, svcUser, "http://localhost:8080")
	require.NoError(t, err)

	return h
}
