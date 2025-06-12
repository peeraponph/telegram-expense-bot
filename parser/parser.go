package parser

import (
	"regexp"
	"strconv"
	"strings"
	"telegram-expense-bot/entity"
	"telegram-expense-bot/util"
	"time"
)

// ParsedMessage struct to hold the parsed message details
func ParseMessage(message string) entity.ExpenseEntry {
	now := util.GetTimestampNow()
	isBackdated := false

	// 1. Rerive tag (if any) from the message e.g. #food, #transport
	tag := ""
	reTag := regexp.MustCompile(`#(\S+)`)
	tagMatch := reTag.FindStringSubmatch(message)
	if len(tagMatch) > 1 {
		tag = "#" + tagMatch[1]
		message = reTag.ReplaceAllString(message, "")
	}

	// 2. Retrieve date from the message e.g. 15/10
	date := now
	reDate := regexp.MustCompile(`\b(\d{1,2})/(\d{1,2})\b`)
	dateMatch := reDate.FindStringSubmatch(message)
	if len(dateMatch) > 2 {
		isBackdated = true             // If we found a date, assume it's backdated
		day := parseInt(dateMatch[1])
		month := parseInt(dateMatch[2])
		year := util.GetBangkokTime().Year()

		nowTime := util.GetBangkokTime()
		dateTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, nowTime.Location())
		date = dateTime.Format("15:04:05 02-01-2006")
		message = reDate.ReplaceAllString(message, "")
	}

	// 3. Retrieve amount from the message e.g. 1,000.50
	reAmount := regexp.MustCompile(`\d+(?:[.,]\d+)?`)

	allMatches := reAmount.FindAllStringIndex(message, -1)

	amount := 0.0
	if len(allMatches) > 0 {
		last := allMatches[len(allMatches)-1]
		raw := message[last[0]:last[1]]
		raw = strings.ReplaceAll(raw, ",", "") // Clear comma
		val, err := strconv.ParseFloat(raw, 64)
		if err == nil {
			amount = val
		}
		message = strings.TrimSpace(message[:last[0]] + message[last[1]:])
	}

	// 4. Check income or expense ? 
	var t string
	if strings.Contains(message, "ขาย") || strings.Contains(message, "ได้") || strings.Contains(message, "รับ") || strings.Contains(message, "เก็บ") {
		t = "รายรับ"
	} else {
		t = "รายจ่าย"
	}

	description := strings.TrimSpace(message)

	// 5. Check if the entry is backdated
	note := "จดเอง"
	if isBackdated {
		note = "จดย้อนหลัง"
	}

	return entity.ExpenseEntry{
		Date:        date,
		Type:        t,
		Description: description,
		Amount:      amount,
		Tag:         tag,
		Note:        note,
	}
}

func parseInt(s string) int {
	n := 0
	for _, ch := range s {
		n = n*10 + int(ch-'0')
	}
	return n
}
