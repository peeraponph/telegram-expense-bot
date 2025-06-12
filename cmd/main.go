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
	err := godotenv.Load(".env.deploy")
	if err != nil {
		log.Println("⚠️ Not Found .env", err)
	}

	// เตรียม Google Sheet Service
	sheetService, err := service.NewGoogleSheetService()
	if err != nil {
		log.Fatal("❌ Google Sheet Service Does Not Initialize:", err)
	}

	// เตรียม Controller
	controller := controller.NewExpenseController(sheetService)

	// สร้าง Telegram Bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("❌  TELEGRAM_BOT_TOKEN Not Found in .env")
	}

	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("❌ Does Not Initialize Telegram Bot", err)
	}
	log.Printf("✅ Initialize Telegram Bot @%s", botAPI.Self.UserName)

	// สร้าง BotHandler
	handler := bot.NewBotHandler(botAPI, controller, sheetService)

	// เริ่มรับ update
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handler.HandleUpdate(update)
		}
	}
}
