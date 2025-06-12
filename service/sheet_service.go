package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"telegram-expense-bot/entity"
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
}

type GoogleSheetService struct {
	srv           *sheets.Service
	spreadsheetID string
}

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

func (s *GoogleSheetService) WriteRow(entry entity.ExpenseEntry) error {
	vr := &sheets.ValueRange{
		Values: [][]interface{}{[]interface{}{entry.Date, entry.Type, entry.Description, entry.Amount, entry.Tag, entry.Note}},
	}

	_, err := s.srv.Spreadsheets.Values.Append(s.spreadsheetID, "Expenses!A:F", vr).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return fmt.Errorf("cannot write row to Google Sheet: %v", err)
	}

	return nil
}

func (s *GoogleSheetService) ReadSheetData() ([]entity.ExpenseEntry, error) {
	resp, err := s.srv.Spreadsheets.Values.Get(s.spreadsheetID, "Expenses!A:F").Do()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from Google Sheet: %v", err)
	}

	var records []entity.ExpenseEntry
	for i, row := range resp.Values {
		if i == 0 {
			continue // ข้าม header
		}

		amount := 0
		fmt.Sscanf(fmt.Sprintf("%v", u.GetSafe(row, 3)), "%d", &amount)

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

// GetTodaySummary สรุปรายรับรายจ่ายของวันนี้
func (s *GoogleSheetService) GetTodaySummary() (string, error) {
	records, err := s.ReadSheetData()
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

func (s *GoogleSheetService) GetMonthSummary() (string, error) {
	records, err := s.ReadSheetData()
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
