package inmemory

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapStorage_GetURL(t *testing.T) {
	type args struct {
		shortURL string
	}
	tests := []struct {
		name    string
		mp      MapStorage
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty shortUrl",
			mp: MapStorage{
				urlMap: map[string]string{
					"dqwdwqd": "https://practicum.yandex.ru/",
				},
				mutex: &sync.RWMutex{},
			},
			args:    args{shortURL: ""},
			wantErr: true,
		},
		{
			name: "find shortUrl",
			mp: MapStorage{
				urlMap: map[string]string{
					"dqwdwqd": "https://practicum.yandex.ru/",
				},
				mutex: &sync.RWMutex{},
			},
			args:    args{shortURL: "dqwdwqd"},
			want:    "https://practicum.yandex.ru/",
			wantErr: false,
		},
		{
			name: "not find shortUrl",
			mp: MapStorage{
				urlMap: map[string]string{
					"dqwdwqd": "https://practicum.yandex.ru/",
				},
				mutex: &sync.RWMutex{},
			},
			args:    args{shortURL: "111111"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mp.GetURL(tt.args.shortURL)
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
	type args struct {
		shortID string
		srcURL  string
	}
	tests := []struct {
		name    string
		mp      MapStorage
		args    args
		wantErr bool
	}{
		{
			name: "empty URL",
			mp: MapStorage{
				urlMap: map[string]string{
					"dqwdwqd": "https://practicum.yandex.ru/",
				},
				mutex: &sync.RWMutex{},
			},
			args:    args{srcURL: "", shortID: ""},
			wantErr: true,
		},
		{
			name: "add URL",
			mp: MapStorage{
				urlMap: map[string]string{
					"dqwdwqd": "https://practicum.yandex.ru/",
				},
				mutex: &sync.RWMutex{},
			},
			args:    args{srcURL: "https://github.com/", shortID: "dsdsdds"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mp.SaveURL(tt.args.shortID, tt.args.srcURL)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Contains(t, tt.mp.urlMap, got)
				assert.Equal(t, tt.mp.urlMap[got], tt.args.srcURL)

				return
			}
			assert.Error(t, err)
		})
	}
}
