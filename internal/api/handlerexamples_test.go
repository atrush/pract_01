package api

import (
	"fmt"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/storage/infile"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandler_Ping() {
	h, err := initExampleHandler()
	if err != nil {
		log.Fatal(err.Error())
	}

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)

	r := NewRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
}

func ExampleHandler_SaveURLHandler() {
	h, err := initExampleHandler()
	if err != nil {
		log.Fatal(err.Error())
	}

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru/"))
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")

	r := NewRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
}

func ExampleHandler_SaveURLJSONHandler() {
	h, err := initExampleHandler()
	if err != nil {
		log.Fatal(err.Error())
	}

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("{\"url\": \"https://www.google.com/\"}"))
	request.Header.Set("Content-Type", "application/json")

	r := NewRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
}

func ExampleHandler_GetURLHandler() {
	h, err := initExampleHandler()
	if err != nil {
		log.Fatal(err.Error())
	}

	request := httptest.NewRequest(http.MethodGet, "/mEoMYa+7", nil)
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")

	r := NewRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
}

func ExampleHandler_GetUserUrls() {
	h, err := initExampleHandler()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := NewRouter(h, false)
	w := httptest.NewRecorder()

	cookie, err := getAuthCookie(r, w)

	//  save some urls
	for i := 0; i < 5; i++ {
		url := fmt.Sprintf("https://practicum.yandex.ru/%v", uuid.New().String())

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))
		request.AddCookie(cookie)
		request.Header.Set("Content-Type", "text/plain; charset=utf-8")

		r.ServeHTTP(w, request)
	}

	//  get user urls list
	request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	request.AddCookie(cookie)
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, request)
}

func ExampleHandler_SaveBatch() {
	h, err := initExampleHandler()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := NewRouter(h, false)
	w := httptest.NewRecorder()

	jsBatch := "[{\"correlation_id\":\"1b8c89ab-608f-4f41-a0fe-e30567e9040a\",\"original_url\":\"http://uw5y0ltlx.ru\"},{\"correlation_id\":\"b3419e03-6055-4055-8ac6-5027b83bb809\",\"original_url\":\"http://vjfl1qozfyl.biz/mfnao4blefvn/ypfsjjlgmu3\"}]"

	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(jsBatch))
	request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, request)
}

func initExampleHandler() (*Handler, error) {
	tstSt, err := infile.NewFileStorage("")
	if err != nil {
		return nil, err
	}

	svcSht, err := service.NewShortURLService(tstSt)
	if err != nil {
		return nil, err
	}

	svcUser, err := service.NewUserService(tstSt)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(svcSht, svcUser, "http://localhost:8080", "")
	if err != nil {
		return nil, err
	}

	return h, nil
}
