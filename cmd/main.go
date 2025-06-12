package main

import (
	"log"
	"os"
	"telegram-expense-bot/bot"
	"telegram-expense-bot/controller"
	"telegram-expense-bot/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {


	_ = godotenv.Load(".env.deploy")


	// Prepair Google Sheet Service
	sheetService, err := service.NewGoogleSheetService()
	if err != nil {
		log.Fatal("❌ Google Sheet Service Does Not Initialize:", err)
	}

	// Prepair Controller
	controller := controller.NewExpenseController(sheetService)

	// Create Telegram Bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("❌  TELEGRAM_BOT_TOKEN Not Found in .env")
	}

	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("❌ Does Not Initialize Telegram Bot", err)
	}
	log.Printf("✅ Initialize Telegram Bot @%s", botAPI.Self.UserName)

	// Create BotHandler
	handler := bot.NewBotHandler(botAPI, controller, sheetService)

	// Run update
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handler.HandleUpdate(update)
		}
	}
}
