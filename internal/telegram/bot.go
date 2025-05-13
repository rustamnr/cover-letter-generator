package telegram

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rustamnr/cover-letter-generator/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	API        *tgbotapi.BotAPI
	APIBaseURL string // URL вашего API
}

// NewBot создает нового Telegram-бота
func NewBot() (*Bot, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true // Включить отладку (опционально)
	log.Printf("Авторизован как %s", bot.Self.UserName)

	apiBaseURL := os.Getenv("API_BASE_URL") // URL вашего API

	return &Bot{API: bot, APIBaseURL: apiBaseURL}, nil
}

// Start запускает обработку сообщений
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Обработка текстовых сообщений
			// Обработка команды /start
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! Выберите действие:")

				// Создаем inline-кнопки
				buttonResume := tgbotapi.NewInlineKeyboardButtonData("Авторизоваться в hh.ru", "hh_login")
				buttonVacancy := tgbotapi.NewInlineKeyboardButtonData("Получить вакансию", "get_vacancy")

				// Создаем клавиатуру с кнопками
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(buttonResume, buttonVacancy),
				)

				// Привязываем клавиатуру к сообщению
				msg.ReplyMarkup = keyboard

				b.API.Send(msg)
				continue
			}

			if update.Message.Text == "/login" {
				telegramID := update.Message.Chat.ID

				// Генерация ссылки для авторизации
				authURL := fmt.Sprintf("%s/auth?telegram_id=%d", b.APIBaseURL, telegramID)

				// Создаем inline-кнопку с ссылкой
				button := tgbotapi.NewInlineKeyboardButtonURL("Авторизоваться в hh.ru", authURL)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(button),
				)

				// Отправляем сообщение с кнопкой
				msg := tgbotapi.NewMessage(telegramID, "Для авторизации перейдите по ссылке:")
				msg.ReplyMarkup = keyboard
				b.API.Send(msg)
			}

			if update.Message.Text == "/resumes" {
				telegramID := update.Message.Chat.ID
				url := fmt.Sprintf("%s/api/resumes", b.APIBaseURL)

				resp, err := http.Get(url)
				if err != nil || resp.StatusCode != http.StatusOK {
					msg := tgbotapi.NewMessage(telegramID, "Ошибка при получении резюме.")
					b.API.Send(msg)
					continue
				}
				msg := tgbotapi.NewMessage(telegramID, "Список резюме получен.")
				b.API.Send(msg)
				continue
			}

			if update.Message.Text == "/vacancy" {
				telegramID := update.Message.Chat.ID
				url := fmt.Sprintf("%s/vacancies/first?telegram_id=%d", b.APIBaseURL, telegramID)

				resp, err := http.Get(url)
				if err != nil || resp.StatusCode != http.StatusOK {
					msg := tgbotapi.NewMessage(telegramID, "Ошибка при получении вакансии.")
					b.API.Send(msg)
					continue
				}
				msg := tgbotapi.NewMessage(telegramID, "Подходящая вакансия получена.")
				b.API.Send(msg)
				continue
			}

			// Пример: обработка ID резюме
			response := b.handleResumeRequest(update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			b.API.Send(msg)
		}

		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			switch callback.Data {
			case "hh_login":
				url := fmt.Sprintf("%s/auth?telegram_id=%d", b.APIBaseURL, callback.From.ID)
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Перейдите по ссылке:")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Авторизация", url)),
				)
				b.API.Send(msg)
			}
			b.API.Request(tgbotapi.NewCallback(callback.ID, ""))
		}
	}
}

func (b *Bot) hhLogin() {
	resp, err := http.Get(b.APIBaseURL + "/auth")
	if err != nil {
		logger.Errorf("Ошибка при запросе к API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Не удалось получить данные. Проверьте ID.")
		return
	}
	logger.Infof("Успешно авторизован в HH API")
}

// handleResumeRequest обрабатывает запросы для резюме
func (b *Bot) handleResumeRequest(input string) string {

	resp, err := http.Get(b.APIBaseURL + "/resumes" + input)
	if err != nil {
		log.Printf("Ошибка при запросе к API: %v", err)
		return "Произошла ошибка при обработке запроса."
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Не удалось получить данные. Проверьте ID."
	}

	return "Резюме успешно обработано!"
}
