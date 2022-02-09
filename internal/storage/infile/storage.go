package infile

import (
	"fmt"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	shortURLRepo *shortURLRepository
	userRepo     *userRepository
	fileName     string
	cache        *cache
}

// Init new file storage
func NewFileStorage(fileName string) (*Storage, error) {
	st := Storage{
		fileName: fileName,
		cache:    newCache(),
	}

	var err error

	st.shortURLRepo, err = newShortURLRepository(st.cache, st.fileName)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	st.userRepo, err = newUserRepository(st.cache)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	if st.fileName != "" {
		if err := st.initFromFile(); err != nil {
			return nil, fmt.Errorf("ошибка инициализации хранилища: %w", err)
		}
	}

	return &st, nil
}

// Return URL repository
func (s *Storage) URL() storage.URLRepository {
	if s.shortURLRepo != nil {
		return s.shortURLRepo
	}
	return s.shortURLRepo
}

// Return User repository
func (s *Storage) User() storage.UserRepository {
	if s.shortURLRepo != nil {
		return s.userRepo
	}
	return s.userRepo
}

// Read all items from file
func (f *Storage) initFromFile() error {
	fileReader, err := newFileReader(f.fileName)
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	data, err := fileReader.ReadAll()
	defer fileReader.Close()
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	//set URL cache
	f.cache.urlCache = data

	if len(data) > 0 {
		for _, v := range data {
			//set URL index
			if f.shortURLRepo.IsAvailableID(v.ShortID) {
				f.cache.shortURLidx[v.ShortID] = v.ShortID
			}

			//set Users cahe
			if v.UserID != "" && f.userRepo.IsAvailableUserID(v.UserID) {
				if _, ok := f.cache.userCache[v.UserID]; !ok {
					f.cache.userCache[v.UserID] = v.UserID
				}
			}
		}
	}

	return nil
}

// Empty, imitate close function
func (s *Storage) Close() {}
