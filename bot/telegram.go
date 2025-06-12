package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"telegram-expense-bot/controller"
	"telegram-expense-bot/entity"
	"telegram-expense-bot/service"
	"time"

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

	chatID := update.Message.Chat.ID

	if update.Message.Photo != nil {
		// Handle photo message
		photos := update.Message.Photo
		if len(photos) == 0 {
			return
		}

		fileID := photos[len(photos)-1].FileID
		file, _ := b.Bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		url := file.Link(os.Getenv("TELEGRAM_BOT_TOKEN"))

		// Download
		resp, err := http.Get(url)
		if err != nil {
			b.Bot.Send(tgbotapi.NewMessage(chatID, "❌ ไม่สามารถโหลดรูปภาพได้"))
			return
		}
		defer resp.Body.Close()

		tmpPath := "tmp.jpg"
		out, _ := os.Create(tmpPath)
		defer os.Remove(tmpPath)
		io.Copy(out, resp.Body)

		amount, err := service.ExtractAmountFromImage(tmpPath)

		if err != nil {
			b.Bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ OCR ไม่สามารถอ่านจำนวนเงินได้: %v", err)))
		} else {
			entry := entity.ExpenseEntry{
				Date:        time.Now().Format("2006-01-02"),
				Type:        "รายจ่าย",
				Description: "จากภาพ OCR",
				Amount:      amount,
				Tag:         "#OCR",
				Note:        "จากรูปภาพ",
			}
			reply := b.Controller.HandleParsedEntry(entry)
			b.Bot.Send(tgbotapi.NewMessage(chatID, reply))
		}
	}

	// Handle text message
	if update.Message.Text != "" {
		text := update.Message.Text
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
}

func (b *BotHandler) handleExport(chatID int64) (string, error) {
	link := os.Getenv("SPREADSHEET_LINK")
	if link == "" {
		link = "⚠️ ยังไม่ได้ตั้งค่า SPREADSHEET_LINK ใน .env"
	}

	reply := fmt.Sprintf("📄 ข้อมูล Google Sheet:\n%s", link)

	// ส่วนสำหรับแนบไฟล์ Excel ไว้เปิดใช้งานภายหลัง
	// excelFile := "export.xlsx"
	// if err := b.Sheet.ExportToExcel(excelFile); err != nil {
	// 	return "", fmt.Errorf("สร้างไฟล์ Excel ไม่สำเร็จ: %v", err)
	// }
	// doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(excelFile))
	// doc.Caption = "📦 ส่งออกข้อมูลรายรับรายจ่าย"
	// if _, err := b.Bot.Send(doc); err != nil {
	// 	return "", fmt.Errorf("ส่งไฟล์ไม่สำเร็จ: %v", err)
	// }
	// _ = os.Remove(excelFile)

	return reply, nil
}
