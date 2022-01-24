package infile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/atrush/pract_01.git/internal/storage"
)

type fileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewFileWriter(filename string) (*fileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации записи в файл: %w", err)
	}

	return &fileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (f *fileWriter) Close() error {
	return f.file.Close()
}

func (f *fileWriter) WriteURL(shortID string, srcURL string) error {
	writeURL := storage.ShortURL{
		ShortID: shortID,
		URL:     srcURL,
	}

	jsURL, err := json.Marshal(writeURL)
	if err != nil {
		return fmt.Errorf("ошибка обработки данных для записи в файл: %w", err)
	}

	if _, err := f.writer.Write(jsURL); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	if err := f.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return f.writer.Flush()
}
