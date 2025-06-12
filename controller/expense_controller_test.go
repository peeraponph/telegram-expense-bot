package controller_test

import (
	"telegram-expense-bot/controller"
	"telegram-expense-bot/entity"
	"testing"
)

// MockSheet implements service.SheetWriter interface
type MockSheet struct {
	Saved []entity.ExpenseEntry
}

func (m *MockSheet) WriteRow(e entity.ExpenseEntry) error {
	m.Saved = append(m.Saved, e)
	return nil
}
func (m *MockSheet) ReadSheetData() ([]entity.ExpenseEntry, error) { return nil, nil }
func (m *MockSheet) GetTodaySummary() (string, error)              { return "", nil }
func (m *MockSheet) GetMonthSummary() (string, error)              { return "", nil }
func (m *MockSheet) ExportToExcel(filename string) error           { return nil }
func (m *MockSheet) AppendToSheet(amount int, source string) error {
	return nil
}

func TestHandleMessage(t *testing.T) {
	mock := &MockSheet{}
	ctrl := controller.NewExpenseController(mock)

	input := "ข้าวมันไก่ 50 #อาหาร"
	reply := ctrl.HandleMessage(input)

	if reply == "" {
		t.Errorf("Expected reply message, got empty string")
	}
	if len(mock.Saved) != 1 {
		t.Errorf("Expected 1 entry saved, got %d", len(mock.Saved))
	}
}
