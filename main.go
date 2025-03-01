package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-bot/handlers"
	"log"
	"os"
)

func main() {
	// Подгружаем енвы
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	// Получаем токен и переменных окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Missing TELEGRAM_BOT_TOKEN")
	}

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = os.Getenv("DEBUG") != ""

	// Разрешаем боту читать обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			go handlers.HandleMessage(bot, update.Message)
		}
	}
}
