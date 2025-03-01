package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// DownloadFileToBuffer скачивает видео и возвращает его в виде буфера
func DownloadFileToBuffer(url string) (*bytes.Buffer, error) {
	// Запрос к файлу
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("❌ Ошибка загрузки видео: %v", err)
	}
	defer resp.Body.Close()

	// Читаем содержимое в буфер
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("❌ Ошибка сохранения в буфер: %v", err)
	}

	return buf, nil
}
