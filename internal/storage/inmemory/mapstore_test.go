package inmemory

import (
	"testing"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapStorage_GetURL(t *testing.T) {

	tests := []struct {
		name         string
		shortID      string
		want         string
		wantErr      bool
		initFixtures func(storage *MapStorage)
	}{
		{
			name:    "empty shortID",
			shortID: "",
			wantErr: true,
			initFixtures: func(storage *MapStorage) {
				storage.SaveURL("dsfdfd", "https://practicum.yandex.ru/")
			},
		},
		{
			name:    "exist shortID",
			shortID: "dsfdfd",
			want:    "https://practicum.yandex.ru/",
			wantErr: false,
			initFixtures: func(storage *MapStorage) {
				storage.SaveURL("dsfdfd", "https://practicum.yandex.ru/")
			},
		},
		{
			name:    "not exist shortID in empty storage",
			want:    "",
			shortID: "dsfdfd",
			wantErr: false,
		},
		{
			name:    "not exist shortID in not empty storage",
			want:    "",
			shortID: "dsfdfd",
			wantErr: false,
			initFixtures: func(storage *MapStorage) {
				storage.SaveURL("dhygff", "https://practicum.yandex.ru/")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewStorage()
			if tt.initFixtures != nil {
				tt.initFixtures(db)
			}
			got, err := db.GetURL(tt.shortID)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)

				return
			}
			assert.Error(t, err)
		})
	}
}

func TestMapStorage_SaveURL(t *testing.T) {
	tests := []struct {
		name         string
		url          storage.ShortURL
		want         string
		wantErr      bool
		initFixtures func(storage *MapStorage)
	}{
		{
			name:    "empty shortID",
			url:     storage.ShortURL{ShortID: "", URL: "https://practicum.yandex.ru/"},
			wantErr: true,
		},
		{
			name:    "empty srcURL",
			url:     storage.ShortURL{ShortID: "xfdafds", URL: ""},
			wantErr: true,
		},
		{
			name:    "add URL",
			url:     storage.ShortURL{ShortID: "xfdafds", URL: "https://practicum.yandex.ru/"},
			wantErr: false,
			want:    "xfdafds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewStorage()
			if tt.initFixtures != nil {
				tt.initFixtures(db)
			}
			got, err := db.SaveURL(tt.url.ShortID, tt.url.URL)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)

				return
			}
			assert.Error(t, err)
		})
	}
}
