package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortLink(t *testing.T) {
	type args struct {
		srcLink string
		salt    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "long link",
			args: args{
				srcLink: "https://www.google.com/search?q=%D1%82%D0%B5%D1%81%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5+golang&oq=%D1%82%D0%B5%D1%81%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5+golang&aqs=chrome..69i57j0i22i30l9.6959j0j15&sourceid=chrome&ie=UTF-8",
				salt:    "3",
			},
		},
		{
			name: "short link",
			args: args{
				srcLink: "https://github.com/",
				salt:    "8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := GenerateShortLink(tt.args.srcLink, tt.args.salt)

			assert.False(t, len(v) == 0 || len(v) > 10, "GenerateShortLink() = %v, длинна генерируемой ссылки должа быть >0 и <10 символов", v)

			regexPatternBase64 := "^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$"
			ok, _ := regexp.MatchString(regexPatternBase64, v)
			assert.True(t, ok, "GenerateShortLink() = %v, генерируемая ссылка должна быть в кодировке BASE64", v)
		})
	}
}
