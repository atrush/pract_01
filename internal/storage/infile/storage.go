package infile

import (
	"errors"
	"fmt"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ storage.Storage = (*Storage)(nil)

//  Storage implements Storage interface, provides storing data in memory and duplicates it to file.
type Storage struct {
	shortURLRepo *shortURLRepository
	userRepo     *userRepository
	fileName     string
	cache        *cache
}

//  NewFileStorage inits new file storage, reads all records from file to memory.
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

//  URL returns urls repository.
func (s *Storage) URL() storage.URLRepository {
	if s.shortURLRepo != nil {
		return s.shortURLRepo
	}
	return s.shortURLRepo
}

//  User returns users repository.
func (s *Storage) User() storage.UserRepository {
	if s.shortURLRepo != nil {
		return s.userRepo
	}
	return s.userRepo
}

//  Ping checks storage connection.
//  Always return error, becous storage database not initialised.
func (s *Storage) Ping() error {
	return errors.New("db not initialized")
}

//  Close is empty, imitates close function
func (s *Storage) Close() {}

//  initFromFile read all items from file to memory.
func (s *Storage) initFromFile() error {
	fileReader, err := newFileReader(s.fileName)
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	data, err := fileReader.ReadAll()
	defer fileReader.Close()
	if err != nil {
		return fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	//  set URL cache
	s.cache.urlCache = data

	if len(data) > 0 {
		for _, v := range data {
			//  set URL index
			existShortID, _ := s.shortURLRepo.Exist(v.ShortID)
			if !existShortID {
				s.cache.shortURLidx[v.ShortID] = v.ID
			}

			//  set Users cahe
			existUser, _ := s.userRepo.Exist(v.UserID)
			if v.UserID != uuid.Nil && !existUser {
				s.cache.userCache[v.UserID] = v.UserID
			}

			//  set srcURL cache
			existURL, _ := s.shortURLRepo.Exist(v.URL)
			if !existURL {
				s.cache.srcURLidx[v.URL] = v.ID
			}
		}
	}

	return nil
}
