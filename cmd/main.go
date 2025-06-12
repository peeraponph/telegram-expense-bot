package main

import (
	"log"
	"net/http"
	"os"
	"telegram-expense-bot/bot"
	"telegram-expense-bot/controller"
	"telegram-expense-bot/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env.deploy file
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

	// Get webhook URL from environment variable
	webhookURL := os.Getenv("WEBHOOK_URL")

	// Initialize Telegram Bot API
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("❌ Does Not Initialize Telegram Bot", err)
	}


	// Set webhook if WEBHOOK_URL is provided
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatal("❌ Failed to create webhook config:", err)
	}
	_, err = botAPI.Request(webhookConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP server to listen for webhook requests
	updates := botAPI.ListenForWebhook("/telegram-webhook")
	go http.ListenAndServe(":8080", nil) // Must Listen on port 8080

	log.Printf("✅ Initialize Telegram Bot @%s", botAPI.Self.UserName)

	// Create BotHandler
	handler := bot.NewBotHandler(botAPI, controller, sheetService)

	for update := range updates {
		if update.Message != nil {
			handler.HandleUpdate(update)
		}
	}
}
