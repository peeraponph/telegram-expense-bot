package sheet

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func WriteRow(row []interface{}) error {
	ctx := context.Background()

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		return fmt.Errorf("ไม่พบค่า SPREADSHEET_ID")
	}

	// โหลด credentials.json
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("อ่าน credentials.json ไม่ได้: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return fmt.Errorf("เชื่อม Google Sheets ไม่ได้: %v", err)
	}

	vr := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "Expenses!A:F", vr).
		ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return fmt.Errorf("เขียนข้อมูลลงชีทไม่สำเร็จ: %v", err)
	}

	return nil
}

type Record struct {
	Date        string
	Type        string
	Description string
	Amount      int
	Tag         string
	Note        string
}

func ReadSheetData() ([]Record, error) {
	ctx := context.Background()
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		return nil, fmt.Errorf("SPREADSHEET_ID ไม่ถูกตั้งค่า")
	}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("อ่าน credentials.json ไม่ได้: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return nil, fmt.Errorf("เชื่อม Google Sheets ไม่สำเร็จ: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, "Expenses!A:F").Do()
	if err != nil {
		return nil, fmt.Errorf("อ่านข้อมูลจาก Google Sheet ไม่ได้: %v", err)
	}

	var records []Record
	for i, row := range resp.Values {
		if i == 0 {
			continue // ข้าม header (ถ้ามี)
		}

		amount := 0
		fmt.Sscanf(fmt.Sprintf("%v", row[3]), "%d", &amount)

		record := Record{
			Date:        getSafe(row, 0),
			Type:        getSafe(row, 1),
			Description: getSafe(row, 2),
			Amount:      amount,
			Tag:         getSafe(row, 4),
			Note:        getSafe(row, 5),
		}

		records = append(records, record)
	}
	return records, nil
}

// GetTodaySummary สรุปรายรับรายจ่ายของวันนี้
func GetTodaySummary() (string, error) {
	records, err := ReadSheetData()
	if err != nil {
		return "", err
	}

	today := time.Now().Format("2006-01-02")
	income := 0
	expense := 0

	for _, rec := range records {
		if rec.Date == today {
			if rec.Type == "รายรับ" {
				income += rec.Amount
			} else if rec.Type == "รายจ่าย" {
				expense += rec.Amount
			}
		}
	}

	balance := income - expense
	return fmt.Sprintf(
		"📊 สรุปวันนี้ (%s):\nรายรับ: %d บาท\nรายจ่าย: %d บาท\nคงเหลือ: %+d บาท",
		today, income, expense, balance,
	), nil
}

// getSafe returns the value at index i from row if it exists, otherwise returns an empty string.
func getSafe(row []interface{}, i int) string {
	if len(row) > i {
		return fmt.Sprintf("%v", row[i])
	}
	return ""
}

func GetMonthSummary() (string, error) {
	records, err := ReadSheetData()
	if err != nil {
		return "", err
	}

	now := time.Now()
	currentMonth := now.Format("2006-01") // เช่น "2025-06"

	income := 0
	expense := 0

	for _, rec := range records {
		if strings.HasPrefix(rec.Date, currentMonth) {
			if rec.Type == "รายรับ" {
				income += rec.Amount
			} else if rec.Type == "รายจ่าย" {
				expense += rec.Amount
			}
		}
	}

	balance := income - expense
	return fmt.Sprintf(
		"📊 สรุปเดือนนี้ (%s):\nรายรับ: %d บาท\nรายจ่าย: %d บาท\nคงเหลือ: %+d บาท",
		now.Format("January 2006"), income, expense, balance,
	), nil
}

func ExportToExcel(filename string) error {
	records, err := ReadSheetData()
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	sheetName := "Expenses"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("สร้างชีตใหม่ไม่สำเร็จ: %v", err)
	}
	f.SetActiveSheet(index)

	headers := []string{"วันที่", "ประเภท", "รายละเอียด", "จำนวนเงิน", "Tag", "หมายเหตุ"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, rec := range records {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), rec.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), rec.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), rec.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), rec.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), rec.Tag)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), rec.Note)
	}

	return f.SaveAs(filename)
}
