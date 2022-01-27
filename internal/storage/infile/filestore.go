package infile

import (
	"errors"
	"fmt"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLStorer = (*FileStorage)(nil)

type FileStorage struct {
	fileName string
	cache    *map[string]string
	mutex    sync.RWMutex
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	fileStorage := FileStorage{
		fileName: fileName,
		mutex:    sync.RWMutex{},
	}
	if err := fileStorage.initFromFile(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	return &fileStorage, nil
}

func (f *FileStorage) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}

	f.mutex.RLock()
	longURL, ok := (*f.cache)[shortID]
	defer f.mutex.RUnlock()

	if ok {
		return longURL, nil
	}

	return "", nil
}

func (f *FileStorage) SaveURL(shortID string, srcURL string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}
	if srcURL == "" {
		return "", errors.New("нельзя сохранить пустое значение")
	}
	if !f.IsAvailableID(shortID) {
		return "", errors.New("id уже существует")
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if err := f.writeToFile(shortID, srcURL); err != nil {
		return "", err
	}

	(*f.cache)[shortID] = srcURL

	return shortID, nil
}
func (f *FileStorage) IsAvailableID(shortID string) bool {
	_, ok := (*f.cache)[shortID]

	return !ok
}

func (f *FileStorage) writeToFile(shortID string, srcURL string) error {
	fileWriter, err := NewFileWriter(f.fileName)
	if err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	defer fileWriter.Close()
	if err := fileWriter.WriteURL(shortID, srcURL); err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	return nil
}

func (f *FileStorage) initFromFile() error {
	fileReader, err := NewFileReader(f.fileName)
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	data, err := fileReader.ReadAll()
	defer fileReader.Close()
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}
	f.cache = data

	return nil
}