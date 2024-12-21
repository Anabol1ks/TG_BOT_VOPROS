package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Создаем нового бота, укажите здесь ваш токен
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Настраиваем обновления
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// ID администраторов (замените на свои)
	admins := []int64{665731091, 1281791401}

	for update := range updates {
		if update.Message == nil { // Игнорируем любые обновления без сообщений
			continue
		}

		switch update.Message.Text {
		case "/start":
			// Отправляем приветствие
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Мы создали этого бота для анонимных вопросов нам.\n\n")
			// Добавляем кнопку
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Задать вопрос"),
				),
			)
			bot.Send(msg)
			time.Sleep(1 * time.Second)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы задаёте нам вопрос, бот отправляет его нам, всё анонимно.\n"+
				"Для того чтобы отправить нам ваш анонимный вопрос, введите команду /ask или нажмите на кнопку ниже 'Задать вопрос'.\n\n")
			bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Если вы вдруг хотите не оставаться анонимами, то можете указать своё имя в конце вопроса.")
			bot.Send(msg)

		case "/ask", "Задать вопрос":
			// Просим пользователя ввести вопрос
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш вопрос, и мы постараемся на него ответить!")
			bot.Send(msg)

		default:
			if update.Message.Voice != nil {
				// Обрабатываем голосовое сообщение
				for _, adminID := range admins {
					voiceMsg := tgbotapi.NewVoice(adminID, tgbotapi.FileID(update.Message.Voice.FileID))
					if _, err := bot.Send(voiceMsg); err != nil {
						log.Printf("Ошибка отправки голосового сообщения администратору %d: %v", adminID, err)
					}
				}
			} else if update.Message.VideoNote != nil {
				// Обрабатываем видеосообщение (кружочек)
				for _, adminID := range admins {
					videoNoteMsg := tgbotapi.NewVideoNote(adminID, update.Message.VideoNote.Length, tgbotapi.FileID(update.Message.VideoNote.FileID))
					if _, err := bot.Send(videoNoteMsg); err != nil {
						log.Printf("Ошибка отправки видеосообщения (кружочка) администратору %d: %v", adminID, err)
					}
				}
			} else {
				// Обрабатываем текстовое сообщение как вопрос
				question := update.Message.Text
				for _, adminID := range admins {
					msg := tgbotapi.NewMessage(adminID, "Анонимный вопрос:\n"+question)
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Ошибка отправки сообщения администратору %d: %v", adminID, err)
					}
				}
			}

			// Подтверждаем отправку
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш вопрос отправлен! Спасибо за участие.")
			bot.Send(msg)
		}
	}
}
