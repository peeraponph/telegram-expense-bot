package service_test

import (
	"telegram-expense-bot/entity"
	"testing"
)

type FakeSheet struct {
	Data []entity.ExpenseEntry
}

func (f *FakeSheet) WriteRow(e entity.ExpenseEntry) error {
	f.Data = append(f.Data, e)
	return nil
}

func (f *FakeSheet) ReadSheetData() ([]entity.ExpenseEntry, error) {
	return f.Data, nil
}

func (f *FakeSheet) GetTodaySummary() (string, error) {
	return "📊 สรุปวันนี้: รายรับ 100 รายจ่าย 50 คงเหลือ 50", nil
}

func (f *FakeSheet) GetMonthSummary() (string, error) {
	return "📊 สรุปเดือนนี้: รายรับ 1000 รายจ่าย 500 คงเหลือ 500", nil
}

func (f *FakeSheet) ExportToExcel(filename string) error {
	return nil
}

func TestGetMonthSummary(t *testing.T) {
	sheet := &FakeSheet{}
	summary, err := sheet.GetMonthSummary()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if summary == "" {
		t.Errorf("Expected summary text, got empty string")
	}
}
