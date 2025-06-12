package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"telegram-expense-bot/entity"
	"telegram-expense-bot/util"
	u "telegram-expense-bot/util"
	"time"

	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetWriter interface {
	WriteRow(entry entity.ExpenseEntry) error
	ReadSheetData() ([]entity.ExpenseEntry, error)
	GetTodaySummary() (string, error)
	GetMonthSummary() (string, error)
	ExportToExcel(filename string) error
	AppendToSheet(amount float64, source string) error
}

type GoogleSheetService struct {
	srv           *sheets.Service
	spreadsheetID string
	b             string
}

// NewGoogleSheetService initializes a new Google Sheet service using credentials from a JSON file.
func NewGoogleSheetService() (*GoogleSheetService, error) {
	ctx := context.Background()

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		return nil, fmt.Errorf("ไม่พบค่า SPREADSHEET_ID")
	}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("cannot read credentials.json file : %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return nil, fmt.Errorf("cannot create sheets service: %v", err)
	}

	return &GoogleSheetService{srv: srv, spreadsheetID: spreadsheetID}, nil
}

// This function writes a single row to the Google Sheet with the provided ExpenseEntry data.
func (s *GoogleSheetService) WriteRow(entry entity.ExpenseEntry) error {
	vr := &sheets.ValueRange{
		Values: [][]interface{}{[]interface{}{entry.Date, entry.Type, entry.Description, entry.Amount, entry.Tag, entry.Note}},
	}

	sheetName := "Expenses"
	writeRange := fmt.Sprintf("%s!A:F", sheetName)

	_, err := s.srv.Spreadsheets.Values.Append(s.spreadsheetID, writeRange, vr).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return fmt.Errorf("cannot write row to Google Sheet: %v", err)
	}

	return nil
}

// This function reads data from the Google Sheet and returns a slice of ExpenseEntry.
func (s *GoogleSheetService) ReadSheetData() ([]entity.ExpenseEntry, error) {

	sheetName := "Expenses"
	writeRange := fmt.Sprintf("%s!A:F", sheetName)

	resp, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, writeRange).Do()

	if err != nil {
		return nil, fmt.Errorf("cannot read data from Google Sheet: %v", err)
	}

	var records []entity.ExpenseEntry
	for i, row := range resp.Values {
		if i == 0 {
			continue // ข้าม header
		}

		amount := 0.0
		fmt.Sscanf(fmt.Sprintf("%v", u.GetSafe(row, 3)), "%f", &amount)

		record := entity.ExpenseEntry{
			Date:        u.GetSafe(row, 0),
			Type:        u.GetSafe(row, 1),
			Description: u.GetSafe(row, 2),
			Amount:      amount,
			Tag:         u.GetSafe(row, 4),
			Note:        u.GetSafe(row, 5),
		}
		records = append(records, record)
	}

	return records, nil
}

// This function summarizes the income and expenses for today.
func (s *GoogleSheetService) GetTodaySummary() (string, error) {
	records, err := s.ReadSheetData()
	if err != nil {
		return "", err
	}

	today := util.GetCurrentDate()
	income := 0.0
	expense := 0.0

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
		"📊 สรุปวันนี้ (%s):\nรายรับ: %.2f บาท\nรายจ่าย: %.2f บาท\nคงเหลือ: %.2f บาท",
		today, income, expense, balance,
	), nil
}

// This function summarizes the income and expenses for the current month.
func (s *GoogleSheetService) GetMonthSummary() (string, error) {
	records, err := s.ReadSheetData()
	if err != nil {
		return "", err
	}

	now := time.Now()
	currentMonth := util.GetCurrentMonth() // e.g. "06-2025"

	income := 0.0
	expense := 0.0

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
		"📊 สรุปเดือนนี้ (%s):\nรายรับ: %.2f บาท\nรายจ่าย: %.2f บาท\nคงเหลือ: %.2f บาท",
		now.Format("January 2006"), income, expense, balance,
	), nil
}

// This function is used to append data to the Google Sheet after processing an image or text input.
func (s *GoogleSheetService) AppendToSheet(amount float64, source string) error {
	sheetName := "Expenses"
	writeRange := fmt.Sprintf("%s!A:F", sheetName)

	// Prepare row (ใส่ข้อมูลแบบขั้นต่ำ)
	row := []interface{}{
		util.GetTimestampNow(), // Date (A)
		"รายจ่าย",              // Type (B)
		"จากภาพ ",           // Description (C)
		amount,                 // Amount (D)
		"OCR",                  // Tag (E)
		source,                 // Note (F)
	}

	rb := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	_, err := s.srv.Spreadsheets.Values.Append(s.spreadsheetID, writeRange, rb).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("unable to append data: %v", err)
	}
	return nil
}

// This function exports the data from the Google Sheet to an Excel file.
func (s *GoogleSheetService) ExportToExcel(filename string) error {
	records, err := s.ReadSheetData()
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	sheetName := "Expenses"

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("cannot create new sheet: %v", err)
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
