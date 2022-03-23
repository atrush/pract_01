package infile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/atrush/pract_01.git/internal/storage/schema"
	"github.com/google/uuid"
)

//  fileReader provides data reading from file.
type fileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

//  newFileReader inits new file reader.
func newFileReader(fileName string) (*fileReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации чтения файла: %w", err)
	}

	return &fileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

//  Close closes file.
func (f *fileReader) Close() error {
	return f.file.Close()
}

//  ReadAll reads all items from file.
func (f *fileReader) ReadAll() (map[uuid.UUID]schema.ShortURL, error) {
	data := make(map[uuid.UUID]schema.ShortURL)
	for f.scanner.Scan() {
		lineURL := schema.ShortURL{}
		if err := json.Unmarshal(f.scanner.Bytes(), &lineURL); err != nil {
			return nil, fmt.Errorf("ошибка обработки данных из файла: %w", err)
		}
		data[lineURL.ID] = lineURL
	}

	if err := f.scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return data, nil
}
