package parser

import (
	"regexp"
	"strings"
	"telegram-expense-bot/entity"
	"time"
)

// ParsedMessage struct to hold the parsed message details	
func ParseMessage(message string) entity.ExpenseEntry {
	now := time.Now().Format("2006-01-02")

	// 1. ดึง tag เช่น #ของกิน
	tag := ""
	reTag := regexp.MustCompile(`#(\S+)`)
	tagMatch := reTag.FindStringSubmatch(message)
	if len(tagMatch) > 1 {
		tag = "#" + tagMatch[1]
		message = reTag.ReplaceAllString(message, "")
	}

	// 2. ดึงวันที่ (ถ้ามี) เช่น 12/6
	date := now
	reDate := regexp.MustCompile(`\b(\d{1,2})/(\d{1,2})\b`)
	dateMatch := reDate.FindStringSubmatch(message)
	if len(dateMatch) > 2 {
		day := dateMatch[1]
		month := dateMatch[2]
		year := time.Now().Year()
		date = time.Date(year, time.Month(parseInt(month)), parseInt(day), 0, 0, 0, 0, time.Local).Format("2006-01-02")
		message = reDate.ReplaceAllString(message, "")
	}

	// 3. ดึงจำนวนเงิน (เลขสุดท้าย)
	reAmount := regexp.MustCompile(`\d+`)
	allMatches := reAmount.FindAllStringIndex(message, -1)

	amount := 0
	if len(allMatches) > 0 {
		last := allMatches[len(allMatches)-1]
		amount = parseInt(message[last[0]:last[1]])
		message = strings.TrimSpace(message[:last[0]] + message[last[1]:])
	}

	// 4. ตรวจว่าเป็นรายรับหรือรายจ่าย
	var t string
	if strings.Contains(message, "ขาย") || strings.Contains(message, "ได้") || strings.Contains(message, "รับ") || strings.Contains(message, "เก็บ") {
		t = "รายรับ"
	} else {
		t = "รายจ่าย"
	}

	description := strings.TrimSpace(message)

	return entity.ExpenseEntry{
		Date:        date,
		Type:        t,
		Description: description,
		Amount:      amount,
		Tag:         tag,
		Note:        "",
	}
}

func parseInt(s string) int {
	n := 0
	for _, ch := range s {
		n = n*10 + int(ch-'0')
	}
	return n
}
