package bot

import (
	"fmt"
	"log"
	"os"
	"telegram-expense-bot/controller"
	"telegram-expense-bot/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	Bot        *tgbotapi.BotAPI
	Controller *controller.ExpenseController
	Sheet      service.SheetWriter
}

func NewBotHandler(bot *tgbotapi.BotAPI, controller *controller.ExpenseController, sheet service.SheetWriter) *BotHandler {
	return &BotHandler{
		Bot:        bot,
		Controller: controller,
		Sheet:      sheet,
	}
}

func (b *BotHandler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	text := update.Message.Text
	chatID := update.Message.Chat.ID

	var reply string
	var err error

	switch text {
	case "/summary":
		reply, err = b.Sheet.GetTodaySummary()

	case "/month":
		reply, err = b.Sheet.GetMonthSummary()

	case "/export":
		reply, err = b.handleExport(chatID)

	default:
		reply = b.Controller.HandleMessage(text)
	}

	if err != nil {
		reply = "❌ " + err.Error()
	}

	msg := tgbotapi.NewMessage(chatID, reply)
	if _, err := b.Bot.Send(msg); err != nil {
		log.Println("❌ ส่งข้อความล้มเหลว:", err)
	}
}

func (b *BotHandler) handleExport(chatID int64) (string, error) {
	link := os.Getenv("SPREADSHEET_LINK")
	if link == "" {
		link = "⚠️ ยังไม่ได้ตั้งค่า SPREADSHEET_LINK ใน .env"
	}

	reply := fmt.Sprintf("📄 ข้อมูล Google Sheet:\n%s", link)

	// // export .xlsx
	// excelFile := "export.xlsx"
	// if err := b.Sheet.ExportToExcel(excelFile); err != nil {
	// 	return "", fmt.Errorf("สร้างไฟล์ Excel ไม่สำเร็จ: %v", err)
	// }

	// doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(excelFile))
	// doc.Caption = "📦 ส่งออกข้อมูลรายรับรายจ่าย"
	// if _, err := b.Bot.Send(doc); err != nil {
	// 	return "", fmt.Errorf("ส่งไฟล์ไม่สำเร็จ: %v", err)
	// }

	// _ = os.Remove(excelFile) // clean up (ไม่บังคับ)
	return reply, nil
}
