package psql

import (
	"strconv"
	"testing"

	"github.com/google/uuid"

	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoSaveAndCheckUser(t *testing.T) {

	cfg, err := pkg.NewConfig()
	require.NoError(t, err)
	if cfg.TestBase == "" {
		return
	}

	tstSt, err := NewTestStorage()
	require.NoError(t, err)

	defer tstSt.Close()
	defer tstSt.DropAll()

	user := st.User{
		ID: uuid.New(),
	}

	require.NoError(t, tstSt.userRepo.AddUser(&user))

	exist, err := tstSt.userRepo.Exist(user.ID)
	require.NoError(t, err, "ошибка при проверке наличия пользователя в базе")
	require.True(t, exist, "пользователь не найден")
}

func TestRepoSaveAndGetURL(t *testing.T) {
	cfg, err := pkg.NewConfig()
	require.NoError(t, err)
	if cfg.TestBase == "" {
		return
	}

	tstSt, err := NewTestStorage()
	require.NoError(t, err)

	defer tstSt.Close()
	defer tstSt.DropAll()

	user := st.User{
		ID: uuid.New(),
	}
	sht := st.ShortURL{
		ID:      uuid.New(),
		ShortID: "111111124",
		URL:     "https://github.com/",
		UserID:  user.ID,
	}

	require.NoError(t, tstSt.userRepo.AddUser(&user))

	require.NoError(t, tstSt.shortURLRepo.SaveURL(&sht))

	exist, err := tstSt.shortURLRepo.Exist(sht.ShortID)
	require.NoError(t, err, "ошибка при проверке наличия shortID в базе")
	require.True(t, exist, "сохраненная ссылка не найдена")

	dbURL, err := tstSt.shortURLRepo.GetURL(sht.ShortID)
	require.NoError(t, err)
	require.Equal(t, sht.URL, dbURL, "сохраненная ссылка %v не соответствует отправленной %v", sht.URL, dbURL)
}

func TestRepoSaveAndGetUserURLarray(t *testing.T) {
	cfg, err := pkg.NewConfig()
	require.NoError(t, err)
	if cfg.TestBase == "" {
		return
	}

	tstSt, err := NewTestStorage()
	require.NoError(t, err)

	defer tstSt.Close()
	defer tstSt.DropAll()

	user := st.User{
		ID: uuid.New(),
	}
	require.NoError(t, tstSt.userRepo.AddUser(&user), "не удалось сохранить пользователя")

	arr, err := tstSt.shortURLRepo.GetUserURLList(user.ID, 100)
	require.NoError(t, err)
	require.Equal(t, len(arr), 0, "для пользователя ссылки не создавались, но получены")

	for i := 0; i < 14; i++ {
		sht := st.ShortURL{
			ID:      uuid.New(),
			ShortID: "11114" + strconv.Itoa(i),
			URL:     "https://github.com/" + strconv.Itoa(i),
			UserID:  user.ID,
		}
		arr = append(arr, sht)
		assert.NoError(t, tstSt.shortURLRepo.SaveURL(&sht), "не удалось сохранить ссылку")
	}

	dbURLs, err := tstSt.shortURLRepo.GetUserURLList(user.ID, 100)
	assert.NoError(t, err)

	found := 0
	for _, dbEl := range dbURLs {
		for _, arrEl := range arr {
			if dbEl.ID == arrEl.ID {
				assert.Equal(t, dbEl, arrEl, "сохраненная ссылка %v не соответствует полученной %v", dbEl, arrEl)
				found++
				continue
			}
		}
	}
	require.Equal(t, found, len(arr), "в полученных ссылках найдено %v ссылок", len(arr)-found)
}
