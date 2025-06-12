package main

import (
	"log"
	"os"

	"telegram-expense-bot/parser"
	"telegram-expense-bot/sheet"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ ไม่พบ .env หรือโหลดไม่ได้")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("❌ กรุณาใส่ TELEGRAM_BOT_TOKEN ใน .env")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("❌ เชื่อมต่อ Telegram Bot ไม่สำเร็จ:", err)
	}

	bot.Debug = true
	log.Printf("✅ เริ่มต้น Telegram Bot @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		text := update.Message.Text
		parsed := parser.ParseMessage(text)

		row := []interface{}{
			parsed.Date,
			parsed.Type,
			parsed.Description,
			parsed.Amount,
			parsed.Tag,
			parsed.Note,
		}

		err := sheet.WriteRow(row)
		var reply string
		if err != nil {
			reply = "❌ เกิดข้อผิดพลาด: " + err.Error()
		} else {
			reply = "✅ บันทึกเรียบร้อย: " + parsed.Description
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
