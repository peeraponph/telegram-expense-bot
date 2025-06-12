package main

import (
	"fmt"
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
		switch update.Message.Text {
		// Summary for today
		case "/summary":
			summary, err := sheet.GetTodaySummary()
			if err != nil {
				summary = "❌ ดึงข้อมูลล้มเหลว: " + err.Error()
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, summary)
			bot.Send(msg)

		// Summary for this month
		case "/month":
			monthSummary, err := sheet.GetMonthSummary()
			if err != nil {
				monthSummary = "❌ เกิดข้อผิดพลาด: " + err.Error()
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, monthSummary)
			bot.Send(msg)


		// Export data to Google Sheets
		case "/export":
			link := os.Getenv("SPREADSHEET_LINK")
			if link == "" {
				link = "⚠️ ยังไม่ได้ตั้งค่า SPREADSHEET_LINK ใน .env"
			}

			reply := fmt.Sprintf("📄 ข้อมูลสรุปรายเดือน Google Sheet:\n%s", link)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))

			// สร้างไฟล์ .xlsx
			// excelFile := "export.xlsx"
			// err := sheet.ExportToExcel(excelFile)
			// if err != nil {
			// 	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ สร้างไฟล์ Excel ไม่สำเร็จ: "+err.Error()))
			// 	break
			// }

			// doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(excelFile))
			// doc.Caption = "📦 ส่งออกข้อมูลรายรับรายจ่าย"
			// bot.Send(doc)

			// // ลบไฟล์ออกจากเครื่อง (optional)
			// os.Remove(excelFile)

		default:
			parsed := parser.ParseMessage(update.Message.Text)

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
}
