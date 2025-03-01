package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

// DownloadInstagramVideo загружает видео через эмулированный браузер
func DownloadInstagramVideo(instagramURL string) (string, error) {
	// Запускаем браузер (headless)
	browser := rod.New().MustConnect()
	defer browser.Close()

	// Включаем антибот-защиту (stealth)
	page := stealth.MustPage(browser)

	// Переходим на Snapsave
	err := page.Navigate("https://snapsave.app/id")
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки страницы: %v", err)
	}

	// Ждем загрузки формы
	page.Timeout(10 * time.Second).MustWaitLoad()

	// Вставляем URL и отправляем форму
	page.MustElement("input[name='url']").MustInput(instagramURL)
	page.MustElement("button[type='submit']").MustClick()

	// Ждем обработки (до 10 сек)
	page.Timeout(10 * time.Second).MustWaitLoad()

	// Ищем ссылку на скачивание видео
	downloadLink, err := page.MustElement("div.download-items__btn a").Attribute("href")
	if err != nil {
		return "", errors.New("не удалось найти ссылку на скачивание видео")
	}

	log.Println("Видео найдено:", *downloadLink)
	return *downloadLink, nil
}
