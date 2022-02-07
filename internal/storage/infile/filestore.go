package infile

import (
	"errors"
	"fmt"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLStorer = (*FileStorage)(nil)
var _ storage.UserStorer = (*FileStorage)(nil)

type FileStorage struct {
	fileName  string
	urlCache  map[string]storage.ShortURL
	userCache map[string]string
	sync.RWMutex
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	fileStorage := FileStorage{
		fileName: fileName,
	}
	if err := fileStorage.initFromFile(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	return &fileStorage, nil
}

func (mp *FileStorage) AddUser(userID string) error {
	if userID == "" {
		return errors.New("нельзя использовать пустой id")
	}
	if !mp.IsAvailableID(userID) {
		return errors.New("id уже существует")
	}

	mp.Lock()
	mp.userCache[userID] = userID
	defer mp.Unlock()

	return nil
}

func (mp *FileStorage) IsAvailableUserID(userID string) bool {
	mp.RLock()
	defer mp.RUnlock()

	_, ok := mp.userCache[userID]

	return !ok
}

func (f *FileStorage) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}

	f.RLock()
	defer f.RUnlock()

	item, ok := (f.urlCache)[shortID]
	if ok {
		return item.URL, nil
	}

	return "", nil
}

func (f *FileStorage) SaveURL(shortID string, srcURL string, userID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}
	if srcURL == "" {
		return "", errors.New("нельзя сохранить пустое значение")
	}
	if !f.IsAvailableID(shortID) {
		return "", errors.New("id уже существует")
	}
	if userID != "" && f.IsAvailableUserID(userID) {
		return "", errors.New("пользователь не найден")
	}

	f.Lock()
	defer f.Unlock()

	if err := f.writeToFile(shortID, srcURL, userID); err != nil {
		return "", err
	}

	f.urlCache[shortID] = storage.ShortURL{
		ShortID: shortID,
		URL:     srcURL,
		UserID:  userID,
	}

	return shortID, nil
}
func (f *FileStorage) IsAvailableID(shortID string) bool {
	f.RLock()
	defer f.RUnlock()

	_, ok := f.urlCache[shortID]

	return !ok
}

func (f *FileStorage) writeToFile(shortID string, srcURL string, userID string) error {
	fileWriter, err := newFileWriter(f.fileName)
	if err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	defer fileWriter.Close()
	if err := fileWriter.WriteURL(shortID, srcURL, userID); err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	return nil
}

func (f *FileStorage) initFromFile() error {
	fileReader, err := newFileReader(f.fileName)
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	data, err := fileReader.ReadAll()
	defer fileReader.Close()
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}
	f.urlCache = data

	f.userCache = map[string]string{}
	if len(data) > 0 {
		for _, v := range data {
			if v.UserID != "" && f.IsAvailableUserID(v.UserID) {
				f.userCache[v.UserID] = v.UserID
			}
		}
	}

	return nil
}
