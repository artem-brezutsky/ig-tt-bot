package handlers

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-bot/services"
)

// HandleMessage обрабатывает сообщения
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.Chat.IsPrivate() == true {
		return
	}

	// Проверяем, есть ли ссылка
	text := message.Text
	if text == "" {
		return
	}

	// Регулярное выражение для Instagram
	instagramRegex := regexp.MustCompile(`https?:\/\/(?:www\.)?instagram\.com\/[\w@?=./-]+`)

	// Если это ссылка на Instagram
	if instagramRegex.MatchString(text) {
		loadingMsg := tgbotapi.NewMessage(message.Chat.ID, "⏳ Загружаю видео...")
		loadingMsg.ReplyToMessageID = message.MessageID
		sentMessage, _ := bot.Send(loadingMsg)

		// Получаем ссылку от Snapsave
		videoURL, err := services.DownloadInstagramVideo(text)
		if err != nil {
			_, _ = bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, sentMessage.MessageID, "❌ Ошибка загрузки видео"))
			return
		}

		// Загружаем видео в буфер
		videoBuffer, err := services.DownloadFileToBuffer(videoURL)
		if err != nil {
			_, _ = bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, sentMessage.MessageID, "❌ Ошибка скачивания видео"))
			return
		}

		// Логируем успешное скачивание
		log.Println("📂 Видео загружено в буфер")

		// Отправляем видео как файл
		videoFile := tgbotapi.FileBytes{
			Name:  "video.mp4",
			Bytes: videoBuffer.Bytes(),
		}

		videoMsg := tgbotapi.NewVideo(message.Chat.ID, videoFile)
		videoMsg.Caption = "🎥 Видео загружено"
		videoMsg.ReplyToMessageID = message.MessageID
		_, _ = bot.Send(videoMsg)

		// Удаляем сообщение "⏳ Загружаю видео..."
		bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, sentMessage.MessageID))
	}
}
