package controller

import (
	"fmt"
	"telegram-expense-bot/entity"
	"telegram-expense-bot/parser"
	"telegram-expense-bot/service"
)

type ExpenseController struct {
	Sheet service.SheetWriter
}

func NewExpenseController(sheet service.SheetWriter) *ExpenseController {
	return &ExpenseController{Sheet: sheet}
}

func (c *ExpenseController) HandleMessage(input string) string {
	entry := parser.ParseMessage(input)
	err := c.Sheet.WriteRow(entry)
	if err != nil {
		return fmt.Sprintf("❌ บันทึกไม่สำเร็จ: %v", err)
	}
	return fmt.Sprintf("✅ บันทึกแล้ว: %s (%.2f)", entry.Description, entry.Amount)
}

func (c *ExpenseController) HandleParsedEntry(entry entity.ExpenseEntry) string {
	err := c.Sheet.WriteRow(entry)
	if err != nil {
		return fmt.Sprintf("❌ บันทึกไม่สำเร็จ: %v", err)
	}
	return fmt.Sprintf("✅ บันทึกแล้ว: %s (%.2f)", entry.Description, entry.Amount)
}
